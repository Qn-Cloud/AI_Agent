package common

const (
	SortByLastMessageTime  = 1
	SortByLastMessageCount = 2
	SortByCharacter        = 3
)

const (
	Normal  = 1
	Deleted = 2
)

// AI Role
const (
	AI_Role_User      = "user"
	AI_Role_Assistant = "assistant"
	AI_Role_System    = "system"
	AI_Role_Unknown   = "unknown"
)

// AI SSE Event
const (
	AI_SSE_Event_Message    = "message"
	AI_SSE_Event_Error      = "error"
	AI_SSE_Event_Done       = "done"
	AI_SSE_Event_End        = "end"
	AI_SSE_Event_Stream     = "stream"
	AI_SSE_Event_Progress   = "progress"
	AI_SSE_Event_Unknown    = "unknown"
	AI_SSE_Event_Thinking   = "thinking"
	AI_SSE_Event_Stream_End = "stream_end"
)
