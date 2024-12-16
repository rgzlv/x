//! # Cleanroom
//!
//! Cleanroom is a CLI program to manage shell environments.

use std::env;
use std::io;
use std::result;

use thiserror::Error;

pub mod args;
pub mod cmds;
mod debug;
pub mod files;
pub mod macros;
pub mod senv;
pub mod table;

type Result<T> = result::Result<T, Err>;

#[derive(Debug, Error)]
enum Err {
	#[error(transparent)]
	Xdg(#[from] xdg::BaseDirectoriesError),
	#[error(transparent)]
	Cmd(#[from] cmds::Err),
	#[error(transparent)]
	IO(#[from] io::Error),
}

fn main() -> Result<()> {
	match cr_main() {
		Ok(ok) => Ok(ok),
		Err(err) => {
			eprintln!("Error: {err}");
			Err(err)
		}
	}
}

fn cr_main() -> Result<()> {
	let cmd = args::CmdMain::from_parse();
	let dirs = xdg::BaseDirectories::with_prefix(env!("CARGO_PKG_NAME"))?;

	match cmd.sub {
		args::CmdMainSub::New { args: args_new } => {
			if let Err(err) = cmds::cmd_new(&cmd.args, &args_new, &dirs) {
				return Err(Err::Cmd(cmds::Err::New(err)));
			}
		}

		args::CmdMainSub::Use { args: args_use } => {
			if let Err(err) = cmds::cmd_use(&cmd.args, &args_use, &dirs) {
				return Err(Err::Cmd(cmds::Err::Use(err)));
			}
		}

		args::CmdMainSub::Rm { args: args_rm } => {
			if let Err(err) = cmds::cmd_rm(&cmd.args, &args_rm, &dirs) {
				return Err(Err::Cmd(cmds::Err::Rm(err)));
			}
		}

		args::CmdMainSub::Ls { args: args_ls } => {
			if let Err(err) = cmds::cmd_ls(&cmd.args, &args_ls, &dirs) {
				return Err(Err::Cmd(cmds::Err::Ls(err)));
			}
		}
	}
	Ok(())
}
