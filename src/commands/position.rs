use std::str::FromStr;

use tetka::games::interface::MoveType;
use tetka::games::{games::chess, interface::PositionType};
use tetka::uxi::{Bundle, Command, Flag, RunError, error};

use super::Context;

pub fn position() -> Command<Context> {
    Command::new(|bundle: Bundle<Context>| {
        let mut ctx = bundle.lock();

        let has_fen = bundle.is_flag_set("fen");
        let has_startpos = bundle.is_flag_set("startpos");

        if has_fen && has_startpos {
            return error!("only one of 'fen <fen>' or 'startpos' flags expected");
        }

        let fen = if has_startpos {
            chess::Position::STARTPOS.to_string()
        } else if has_fen {
            bundle.get_array_flag("fen").unwrap().join(" ")
        } else {
            return error!("one of 'startpos' or 'fen <fen>' flags expected");
        };

        ctx.position = chess::Position::from_str(&fen)?;

        if bundle.is_flag_set("moves") {
            for mov in bundle.get_array_flag("moves").unwrap() {
                ctx.position = ctx
                    .position
                    .after_move::<false>(chess::Move::from_str(&mov, &ctx.position)?);
            }
        }

        Ok(())
    })
    .flag("fen", Flag::Array(6))
    .flag("startpos", Flag::Boolean)
    .flag("moves", Flag::Variadic)
}
