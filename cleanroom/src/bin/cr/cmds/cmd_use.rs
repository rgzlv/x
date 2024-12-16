use std::io;
use std::process;
use std::result;

use thiserror::Error;

use crate::args;
use crate::debug::{dbgfmt, DebugPanic};
use crate::files;
use crate::senv;
use crate::table;

type Result<T> = result::Result<T, Err>;

#[derive(Debug, Error)]
pub enum Err {
	#[error(transparent)]
	Files(#[from] files::Err),
	#[error(transparent)]
	Table(#[from] table::Err),
	#[error(transparent)]
	IO(#[from] io::Error),
	#[error(transparent)]
	ShellEnv(#[from] senv::Err),

	#[error("Couldn't convert a `PathBuf` directory to `&str`")]
	DirToStr,
}

pub fn cmd_use(
	_args_main: &args::CmdMainArgs,
	args_use: &args::SubCmdUseArgs,
	dirs: &xdg::BaseDirectories,
) -> Result<()> {
	let shell_env = senv::Senv::new_xdg(&args_use.name, dirs)?;
	let env_table = table::Root::from_env(&args_use.name, dirs)?;
	let shell_args = env_table.get_shell_args(&args_use.name, dirs)?;
	dbgfmt!("Using config: {:#?}", env_table);
	dbgfmt!("Calling with args: {:?}", shell_args);

	// Delete `bin` dir and don't return error if it's a "NotFound" error.
	if shell_env.files.bin_dir.try_exists()? {
		match std::fs::remove_dir_all(&shell_env.files.bin_dir) {
			Ok(_) => (),
			Err(err) => {
				if let std::io::ErrorKind::NotFound = err.kind() {
				} else {
					return Err(Err::IO(err)).dp();
				}
			}
		}
	}

	std::fs::create_dir_all(&shell_env.files.bin_dir).dp()?;

	env_table.bin.inherit_bins(&shell_env.files.data_dir)?;

	let mut shell = process::Command::new(env_table.shell.bin);
	let mut shell = shell.args(shell_args).env_clear();

	let shell_env_vars = env_table.vars.to_env()?;
	#[allow(clippy::iter_over_hash_type)]
	for (k, v) in shell_env_vars {
		shell = shell.env(k, v);
	}

	let mut shell_path = env_table
		.bin
		.inherit_dirs
		.iter()
		.map(|dir| dir.to_str())
		.collect::<Option<Vec<_>>>()
		.ok_or(Err::DirToStr)
		.dp()?
		.join(":");
	if !env_table.bin.inherit.is_empty()
		|| !env_table.bin.inherit_rename.is_empty()
		|| env_table.bin.coreutils
	{
		let env_bin_dir_str = shell_env
			.files
			.bin_dir
			.to_str()
			.ok_or(Err::DirToStr)?
			.to_owned();
		shell_path = env_bin_dir_str + ":" + &shell_path;
	}
	shell.env("PATH", shell_path);

	let mut shell = shell.spawn().dp()?;
	shell.wait().dp()?;

	Ok(())
}
