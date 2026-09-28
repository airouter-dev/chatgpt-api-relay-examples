mod artifacts;
mod client;
mod config;
mod error;
mod types;
mod workflow;

use std::path::PathBuf;

use clap::{Parser, Subcommand};

use crate::{
    artifacts::save_image,
    client::ImageClient,
    config::Settings,
    error::Result,
    workflow::Workflow,
};

const DEFAULT_MODEL: &str = "gpt-image-2";
const DEFAULT_SIZE: &str = "1024x1024";

#[derive(Debug, Parser)]
#[command(author, version, about = "AI-ROUTER image generation in Rust")]
struct Cli {
    /// AI-ROUTER-compatible endpoint, for example https://api.ai-router.dev/v1.
    #[arg(long, env = "AI_ROUTER_BASE_URL")]
    base_url: Option<String>,

    #[command(subcommand)]
    command: Command,
}

#[derive(Debug, Subcommand)]
enum Command {
    /// Turn a text prompt into an image.
    Generate {
        #[arg(long)]
        prompt: String,
        #[arg(long, env = "AI_ROUTER_MODEL", default_value = DEFAULT_MODEL)]
        model: String,
        #[arg(long, env = "AI_ROUTER_SIZE", default_value = DEFAULT_SIZE)]
        size: String,
        #[arg(long, default_value = "output/rust-generated.png")]
        output: PathBuf,
    },
    /// Transform an existing image with the /images/edits endpoint.
    Edit {
        #[arg(long)]
        image: PathBuf,
        #[arg(long)]
        prompt: String,
        #[arg(long, env = "AI_ROUTER_MODEL", default_value = DEFAULT_MODEL)]
        model: String,
        #[arg(long, env = "AI_ROUTER_SIZE", default_value = DEFAULT_SIZE)]
        size: String,
        #[arg(long, default_value = "output/rust-edited.png")]
        output: PathBuf,
    },
    /// Run the complete generate -> edit workflow in one command.
    Workflow {
        #[arg(long)]
        prompt: String,
        #[arg(long)]
        edit_prompt: String,
        #[arg(long, env = "AI_ROUTER_MODEL", default_value = DEFAULT_MODEL)]
        model: String,
        #[arg(long, env = "AI_ROUTER_SIZE", default_value = DEFAULT_SIZE)]
        size: String,
        #[arg(long, default_value = "output/rust-generated.png")]
        generated: PathBuf,
        #[arg(long, default_value = "output/rust-edited.png")]
        edited: PathBuf,
    },
}

#[tokio::main]
async fn main() -> Result<()> {
    let cli = Cli::parse();
    let mut settings = Settings::from_env()?;
    if let Some(base_url) = cli.base_url {
        settings.base_url = base_url.trim_end_matches('/').to_owned();
    }
    let client = ImageClient::new(settings)?;

    match cli.command {
        Command::Generate {
            prompt,
            model,
            size,
            output,
        } => {
            let image = client.generate(&prompt, &model, &size).await?;
            let bytes = client.image_bytes(&image).await?;
            save_image(&output, &bytes).await?;
            println!("Generated image saved to {}", output.display());
        }
        Command::Edit {
            image,
            prompt,
            model,
            size,
            output,
        } => {
            let result = client.edit(&image, &prompt, &model, &size).await?;
            let bytes = client.image_bytes(&result).await?;
            save_image(&output, &bytes).await?;
            println!("Edited image saved to {}", output.display());
        }
        Command::Workflow {
            prompt,
            edit_prompt,
            model,
            size,
            generated,
            edited,
        } => {
            Workflow::new(&client)
                .run(&prompt, &edit_prompt, &model, &size, &generated, &edited)
                .await?;
            println!("Generated image saved to {}", generated.display());
            println!("Edited image saved to {}", edited.display());
        }
    }

    Ok(())
}
