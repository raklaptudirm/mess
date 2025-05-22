use std::{
    str::FromStr,
    sync::{Arc, Mutex},
};

use tetka::{
    games::{games::chess, interface::PositionType},
    uxi::Command,
};

use crate::core::Searcher;

pub use bench::*;
pub use go::*;
pub use perft::*;
pub use position::*;

mod bench;
mod go;
mod perft;
mod position;

pub struct Context {
    pub position: chess::Position,
    pub searcher: Arc<Mutex<Searcher>>,
}

impl Default for Context {
    fn default() -> Self {
        let position = chess::Position::from_str(chess::Position::STARTPOS).unwrap();
        Context {
            position,
            searcher: Arc::new(Mutex::new(Searcher::new())),
        }
    }
}

pub fn d() -> Command<Context> {
    Command::new(|bundle| {
        let ctx = bundle.lock();
        println!("{}", ctx.position);

        Ok(())
    })
}

pub fn ucinewgame() -> Command<Context> {
    Command::new(|bundle| {
        let mut ctx = bundle.lock();
        ctx.position = chess::Position::from_str(chess::Position::STARTPOS)?;

        Ok(())
    })
}

pub fn protocol() -> Command<Context> {
    Command::new(|bundle| {
        let ctx = bundle.lock();
        println!("info string current protocol: {}", ctx.protocol());
        Ok(())
    })
}
