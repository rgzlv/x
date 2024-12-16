use std::io;
use std::result;

use thiserror::Error;

use crate::args;
use crate::senv;

type Result<T> = result::Result<T, Err>;

#[derive(Debug, Error)]
pub enum Err {
	#[error(transparent)]
	IO(#[from] io::Error),
	#[error(transparent)]
	ShellEnv(#[from] senv::Err),
}

pub fn cmd_rm(
	_args_main: &args::CmdMainArgs,
	args_rm: &args::SubCmdRmArgs,
	dirs: &xdg::BaseDirectories,
) -> Result<()> {
	senv::Senv::new_xdg(&args_rm.name, dirs)?.rm()?;

	Ok(())
}
