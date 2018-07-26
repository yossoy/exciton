package menu

type MenuRole string

const (
	RoleAbout              MenuRole = "about"
	RoleHide                        = "hide"
	RoleHideOthers                  = "hideothers"
	RoleUnhide                      = "unhide"
	RoleFront                       = "front"
	RoleUndo                        = "undo"
	RoleRedo                        = "redo"
	RoleCut                         = "cut"
	RoleCopy                        = "copy"
	RolePaste                       = "paste"
	RoleDelete                      = "delete"
	RolePasteAndMatchStyle          = "pasteandmatchstyle"
	RoleSelectAll                   = "selectall"
	RoleStartSpeaking               = "startspeaking"
	RoleStopSpeaking                = "stopspeaking"
	RoleMinimize                    = "minimize"
	RoleClose                       = "close"
	RoleZoom                        = "zoom"
	RoleQuit                        = "quit"
	RoleToggleFullscreen            = "togglefullscreen"
	RoleGoBack                      = "back"
	RoleGoForward                   = "forward"

	RoleServices = "services"
	RoleWindow   = "window"
	RoleHelp     = "help"
)

func roleIsMenuedRole(r MenuRole) bool {
	switch r {
	case RoleServices, RoleWindow, RoleHelp:
		return true
	default:
		return false
	}
}
