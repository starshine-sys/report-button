package main

import "time"

// The timeout before a message can be reported again
const timeout = 5 * time.Minute

const reportError = "There was an error sending your report! Try the <:__reportTo_mods:745306800265494598> reaction, or DM <@!596138399589728262>."
const reportOK = "I've forwarded your report!"

// for transparency, could also set this to something different than reportOK
const timeoutMsg = reportOK
