use std::str::FromStr;

use tetka::{
    games::{games::chess, interface::PositionType},
    uxi::Command,
};

use super::Context;

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
        println!("current protocol: {}", ctx.protocol());
        Ok(())
    })
}
