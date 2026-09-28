use std::io;

use thiserror::Error;

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Debug, Error)]
pub enum Error {
    #[error("AI_ROUTER_API_KEY (or OPENAI_API_KEY) is not set")]
    MissingApiKey,
    #[error("HTTP request failed: {0}")]
    Http(#[from] reqwest::Error),
    #[error("API returned HTTP {status}: {message}")]
    Api { status: u16, message: String },
    #[error("could not decode API response: {0}")]
    Json(#[from] serde_json::Error),
    #[error("image response did not include b64_json or url")]
    EmptyImageData,
    #[error("invalid image response: {0}")]
    InvalidResponse(String),
    #[error("file operation failed: {0}")]
    Io(#[from] io::Error),
    #[error("base64 image data is invalid: {0}")]
    Base64(#[from] base64::DecodeError),
}
