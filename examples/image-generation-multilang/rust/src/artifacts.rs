use std::path::Path;

use tokio::fs;

use crate::error::Result;

/// Persist bytes separately from transport concerns so changing API providers
/// does not change the output pipeline.
pub async fn save_image(path: &Path, bytes: &[u8]) -> Result<()> {
    if let Some(parent) = path.parent() {
        if !parent.as_os_str().is_empty() {
            fs::create_dir_all(parent).await?;
        }
    }
    fs::write(path, bytes).await?;
    Ok(())
}
