use std::path::Path;

use base64::{engine::general_purpose::STANDARD, Engine};
use reqwest::{multipart, Client, Response};
use serde::{de::DeserializeOwned, Serialize};
use serde_json::Value;

use crate::{
    config::Settings,
    error::{Error, Result},
    types::{ImageData, ImageResponse},
};

#[derive(Clone)]
pub struct ImageClient {
    http: Client,
    settings: Settings,
}

#[derive(Debug, Serialize)]
struct GenerateRequest<'a> {
    model: &'a str,
    prompt: &'a str,
    size: &'a str,
    n: u8,
    response_format: &'static str,
}

impl ImageClient {
    pub fn new(settings: Settings) -> Result<Self> {
        let http = Client::builder()
            .timeout(settings.timeout)
            .user_agent("ai-router-image-demo/0.1")
            .build()?;
        Ok(Self { http, settings })
    }

    pub async fn generate(&self, prompt: &str, model: &str, size: &str) -> Result<ImageData> {
        if prompt.trim().is_empty() {
            return Err(Error::InvalidResponse("prompt must not be empty".into()));
        }
        let payload = GenerateRequest {
            model,
            prompt,
            size,
            n: 1,
            response_format: "b64_json",
        };
        let response = self
            .http
            .post(self.endpoint("images/generations"))
            .bearer_auth(&self.settings.api_key)
            .json(&payload)
            .send()
            .await?;
        let payload: ImageResponse = self.parse_response(response).await?;
        payload
            .data
            .into_iter()
            .next()
            .ok_or(Error::InvalidResponse("API returned no images".into()))
    }

    pub async fn edit(
        &self,
        image_path: &Path,
        prompt: &str,
        model: &str,
        size: &str,
    ) -> Result<ImageData> {
        if prompt.trim().is_empty() {
            return Err(Error::InvalidResponse("edit prompt must not be empty".into()));
        }
        let image = tokio::fs::read(image_path).await?;
        if image.is_empty() {
            return Err(Error::InvalidResponse("input image is empty".into()));
        }
        let file_name = image_path
            .file_name()
            .and_then(|name| name.to_str())
            .unwrap_or("input.png")
            .to_owned();
        let part = multipart::Part::bytes(image)
            .file_name(file_name)
            .mime_str(image_mime(image_path));
        let part = part.map_err(Error::Http)?;
        let form = multipart::Form::new()
            .text("model", model.to_owned())
            .text("prompt", prompt.to_owned())
            .text("size", size.to_owned())
            .text("response_format", "b64_json")
            .part("image", part);
        let response = self
            .http
            .post(self.endpoint("images/edits"))
            .bearer_auth(&self.settings.api_key)
            .multipart(form)
            .send()
            .await?;
        let payload: ImageResponse = self.parse_response(response).await?;
        payload
            .data
            .into_iter()
            .next()
            .ok_or(Error::InvalidResponse("API returned no edited images".into()))
    }

    /// Convert either b64_json or a provider URL into bytes for the artifact layer.
    pub async fn image_bytes(&self, image: &ImageData) -> Result<Vec<u8>> {
        if let Some(encoded) = image.b64_json.as_deref() {
            let encoded = encoded
                .split_once(',')
                .map(|(_, payload)| payload)
                .unwrap_or(encoded)
                .replace(&['\n', '\r', ' ', '\t'][..], "");
            let bytes = STANDARD.decode(encoded)?;
            if bytes.is_empty() {
                return Err(Error::InvalidResponse("provider returned an empty image".into()));
            }
            return Ok(bytes);
        }
        if let Some(url) = image.url.as_deref() {
            let response = self.http.get(url).send().await?;
            return self.read_bytes(response).await;
        }
        Err(Error::EmptyImageData)
    }

    fn endpoint(&self, path: &str) -> String {
        format!("{}/{}", self.settings.base_url.trim_end_matches('/'), path)
    }

    async fn parse_response<T: DeserializeOwned>(&self, response: Response) -> Result<T> {
        let status = response.status();
        let bytes = response.bytes().await?;
        if !status.is_success() {
            return Err(Error::Api {
                status: status.as_u16(),
                message: api_error_message(&bytes),
            });
        }
        Ok(serde_json::from_slice(&bytes)?)
    }

    async fn read_bytes(&self, response: Response) -> Result<Vec<u8>> {
        let status = response.status();
        let bytes = response.bytes().await?;
        if !status.is_success() {
            return Err(Error::Api {
                status: status.as_u16(),
                message: "image download failed".into(),
            });
        }
        Ok(bytes.to_vec())
    }
}

fn api_error_message(bytes: &[u8]) -> String {
    serde_json::from_slice::<Value>(bytes)
        .ok()
        .and_then(|value| {
            value
                .get("error")
                .and_then(|error| error.get("message"))
                .and_then(Value::as_str)
                .map(str::to_owned)
        })
        .unwrap_or_else(|| String::from_utf8_lossy(bytes).trim().to_owned())
}

fn image_mime(path: &Path) -> &'static str {
    match path
        .extension()
        .and_then(|extension| extension.to_str())
        .map(|extension| extension.to_ascii_lowercase())
        .as_deref()
    {
        Some("jpg") | Some("jpeg") => "image/jpeg",
        Some("webp") => "image/webp",
        Some("gif") => "image/gif",
        _ => "image/png",
    }
}
