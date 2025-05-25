use tetka::{
    games::{
        games::chess,
        interface::{MoveType, PositionType},
    },
    uxi::{Bundle, Command, Flag, RunError, error, lock},
};

use crate::core::Limits;

use super::Context;

pub fn go() -> Command<Context> {
    Command::new(|bundle: Bundle<Context>| {
        lock! {
            bundle > ctx =>
            let position = ctx.position.clone(); // Get the position to search
            let searcher = ctx.searcher.clone(); // Get the previous search state
        }

        match searcher.lock() {
            Ok(mut searcher) => {
                let mut nodes = 0;

                let limits = parse_limits(&bundle, &position)?;

                let bestmove = searcher.search(&position, limits, &mut nodes);

                println!("bestmove {}", bestmove.with_position(&position));

                Ok(())
            }
            Err(_) => error!("unable to acquire searcher"),
        }
    })
    // Flags for reporting the current time situation.
    .flag("binc", Flag::Single)
    .flag("winc", Flag::Single)
    .flag("btime", Flag::Single)
    .flag("wtime", Flag::Single)
    .flag("movestogo", Flag::Single)
    // Flags for imposing other search limits.
    .flag("depth", Flag::Single)
    .flag("nodes", Flag::Single)
    .flag("movetime", Flag::Single)
    // Flags for setting the search type.
    // .flag("ponder", Flag::Single)
    .flag("infinite", Flag::Single)
    // This command should be run in a separate thread so that the Client
    // can still respond to and run other Commands while this one is running.
    .parallelize(true)
}

fn parse_limits(bundle: &Bundle<Context>, position: &chess::Position) -> Result<Limits, RunError> {
    ////////////////////////////////////////////
    // Check which of the limit flags are set //
    ////////////////////////////////////////////

    // Standard Time Control flags
    let binc = bundle.is_flag_set("binc");
    let winc = bundle.is_flag_set("winc");
    let btime = bundle.is_flag_set("btime");
    let wtime = bundle.is_flag_set("wtime");

    // Other limit flags
    let depth = bundle.is_flag_set("depth");
    let nodes = bundle.is_flag_set("nodes");
    let movetime = bundle.is_flag_set("movetime");

    // Infinite flag
    let infinite = bundle.is_flag_set("infinite");

    ///////////////////////////////////////////////////////
    // Ensure that the given flag configuration is valid //
    ///////////////////////////////////////////////////////

    let std_tc = btime || wtime || binc || winc;
    let oth_tc = depth || nodes || movetime;

    // No other time control flags may be set alongside 'infinite'.
    if infinite && (std_tc || oth_tc) {
        return error!("bad flag set: time control flags set alongside infinite");
    }

    ////////////////////////////////////////////
    // Parse the provided search/perft limits //
    ////////////////////////////////////////////

    Ok(Limits {
        maxnodes: bundle.get_parsed_flag("nodes")?,
        maxdepth: bundle.get_parsed_flag("depth")?,
        movetime: if std_tc {
            let (time, incr) = match position.side_to_move() {
                chess::Color::Black => ("btime", "binc"),
                chess::Color::White => ("wtime", "winc"),
            };

            let time: u128 = bundle.get_parsed_flag(time)?.unwrap_or(0);
            let incr: u128 = bundle.get_parsed_flag(incr)?.unwrap_or(0);

            Some((time / 20 + incr / 2).max(1))
        } else {
            bundle.get_parsed_flag("movetime")?
        },
        movestogo: bundle.get_parsed_flag("movestogo")?,
    })
}
