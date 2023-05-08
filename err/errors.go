package err

const (
	ACTION_FORBIDDEN = "action.forbidden" // payload: action name
	BAD_REQUEST      = "request.bad"      // payload: expected action name
	INVALID_STATE    = "state.invalid"    // payload: state name
	PLAYER_NOT_FOUND = "player.not-found" // payload: player name
	CARD_NOT_FOUND   = "card.not-found"   // payload: card id
	CARD_NOT_PLAYED  = "card.not-played"  // payload: card id
)
