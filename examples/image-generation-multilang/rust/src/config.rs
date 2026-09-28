use std::{env, time::Duration};

use crate::error::{Error, Result};

const DEFAULT_BASE_URL: &str = "https://api.ai-router.dev/v1";

#[derive(Clone, Debug)]
pub struct Settings {
    pub base_url: String,
    pub api_key: String,
    pub timeout: Duration,
}

impl Settings {
    pub fn from_env() -> Result<Self> {
        let api_key = env::var("AI_ROUTER_API_KEY")
            .or_else(|_| env::var("OPENAI_API_KEY"))
            .map_err(|_| Error::MissingApiKey)?;
        let api_key = api_key.trim().to_owned();
        if api_key.is_empty() {
            return Err(Error::MissingApiKey);
        }

        let timeout_seconds = env::var("AI_ROUTER_TIMEOUT_SECONDS")
            .ok()
            .and_then(|value| value.parse::<u64>().ok())
            .unwrap_or(180);

        Ok(Self {
            base_url: env::var("AI_ROUTER_BASE_URL")
                .unwrap_or_else(|_| DEFAULT_BASE_URL.to_owned())
                .trim_end_matches('/')
                .to_owned(),
            api_key,
            timeout: Duration::from_secs(timeout_seconds.max(1)),
        })
    }
}
