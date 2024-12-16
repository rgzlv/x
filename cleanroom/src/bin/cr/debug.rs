//! Macros and `impl`s for debugging.

use std::fmt::{Debug, Display};

// Panic on debug builds or return `Result<T, E>` on release builds.
pub trait DebugPanic<T, E>
where
	E: Display + Debug,
{
	// Short for "debug panic"
	fn dp(self) -> Result<T, E>;
}

#[allow(clippy::panic, clippy::panic_in_result_fn)]
impl<T, E> DebugPanic<T, E> for Result<T, E>
where
	E: Display + Debug,
{
	fn dp(self) -> Self {
		if cfg!(debug_assertions) {
			if let Err(err) = self {
				panic!("{err}");
			} else {
				self
			}
		} else {
			self
		}
	}
}

// `dbg!` doesn't use `Display` so on debug builds print using `println!`.
macro_rules! dbgfmt {
	($($exs:expr),+) => {
		if cfg!(debug_assertions) {
			println!($($exs),+);
		}
	};
}
pub(crate) use dbgfmt;
