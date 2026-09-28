use std::path::Path;

use crate::{
    artifacts::save_image,
    client::ImageClient,
    error::Result,
};

/// Application-level composition for the reviewable generate -> edit flow.
/// Transport and filesystem details stay behind their respective modules.
pub struct Workflow<'a> {
    client: &'a ImageClient,
}

impl<'a> Workflow<'a> {
    pub fn new(client: &'a ImageClient) -> Self {
        Self { client }
    }

    pub async fn run(
        &self,
        prompt: &str,
        edit_prompt: &str,
        model: &str,
        size: &str,
        generated: &Path,
        edited: &Path,
    ) -> Result<()> {
        let generated_image = self.client.generate(prompt, model, size).await?;
        let generated_bytes = self.client.image_bytes(&generated_image).await?;
        save_image(generated, &generated_bytes).await?;

        let edited_image = self
            .client
            .edit(generated, edit_prompt, model, size)
            .await?;
        let edited_bytes = self.client.image_bytes(&edited_image).await?;
        save_image(edited, &edited_bytes).await?;
        Ok(())
    }
}
