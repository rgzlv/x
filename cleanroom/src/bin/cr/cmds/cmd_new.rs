use std::result;

use thiserror::Error;

use crate::args;
use crate::files;
use crate::senv;

type Result<T> = result::Result<T, Err>;

#[derive(Debug, Error)]
pub enum Err {
	#[error(transparent)]
	Files(#[from] files::Err),
	#[error(transparent)]
	ShellEnv(#[from] senv::Err),
}

pub fn cmd_new(
	_args_main: &args::CmdMainArgs,
	args_new: &args::SubCmdNewArgs,
	dirs: &xdg::BaseDirectories,
) -> Result<()> {
	senv::Senv::create_new_xdg(&args_new.name, dirs)?;

	Ok(())
}
