package configs

import "time"

//Const config system
const (
	ServerDB = "tlu.cloudapp.net"
	// /neo4jURL = "Bolt://neo4j:tlis2016@tlu.cloudapp.net:7687"
	URLDB   = "http://neo4j:madawg00@" + ServerDB + ":7474/db/data/"
	APIPort = "8080"
)

// Security const
const (
	JWTSecretKey         = "socialnetworkTLU"         // for login token
	JWTSecretKeyGetImage = "socialnetworkTLUgetimage" // for get img from TLU
	JWTTimeExpire        = time.Hour * 720            // time to expire token
)

// Mail Sender const
const (
	MailAddress = "kien.laohac@gmail.com"
	MailKey     = ""
)

// Twilio const
const (
	TwilioSID       = "1"
	TwilioAPISecret = "@"
)

// Sort const
const (
	ILimit = 25
	SLimit = "25"
	ISkip  = 0
	SSkip  = "0"
	SSort  = "-created_at"
)

//ObjectType const
const (
	PostObject    = 1
	UserObject    = 2
	CommentObject = 3
	LikeObject    = 4
	PageObject    = 5
	GroupObject   = 6
)

// Const privacy post
const (
	Public           = 1
	ShareToFollowers = 2
	Private          = 3
)

// Const privacy group
const (
	// All users within the community can join to this group
	IPublicGroup = 1
	SPublicGroup = "public"
	// Only members of the group can see posts to this group, need request to injoin group
	IPrivateGroup = 2
	SPrivateGroup = "private"
)

// const role member in groupID
const (
	IRoleMember  = 1
	SRoleMember  = "member"
	IRoleAdmin   = 2
	SRoleAdmin   = "admin"
	IRoleCreator = 3
	SCreator     = "creator"
	IRoleBlock   = 4
	SRoleBlock   = "block"
)

// Const status group membership request
const (
	IPending    = 1
	SPending    = "pending"
	IMember     = 4
	SMember     = "member"
	IBlocked    = 5
	SBlocked    = "blocked"
	IAdmin      = 6
	SAdmin      = "admin"
	ICanJoin    = 7
	SCanJoin    = "can_join"
	ICanRequest = 8
	SCanRequest = "can_request"
)

// Const type in get users in post/comment
const (
	ICanMention = 1
	SCanMention = "can_mention"
	IMentioned  = 2
	SMentioned  = "mentioned"
)

// Const type in get users in post
const (
	IFollowedPost  = 3
	SFollowedPost  = "followed"
	ICommentedPost = 4
	SCommentedPost = "commented"
	ILikedPost     = 5
	SLikedPost     = "liked"
)

// Const channel
const (
	IIsAdminChannel  = 1
	IFollowedChannel = 2
)

// Const type action notification
const (
	IActionLike          = 1 // last user liked post
	SActionLike          = "like"
	IActionComment       = 2 // last user commented post
	SActionComment       = "comment"
	IActionShare         = 3 // last user shared post
	SActionShare         = "share"
	IActionMention       = 4 // last user mentioned me in a post
	SActionMention       = "mention"
	IActionPost          = 5 // last post created by a user
	SActionPost          = "post"
	IActionFollow        = 6 // last object followed created by a user
	SActionFollow        = "follow"
	IActionPostPhoto     = 7 // last photo created by a user
	SActionPostPhoto     = "photo_post"
	IActionPostStatus    = 8 // last status created by a user
	SActionPostStatus    = "status_post"
	IActionLikedPost     = 9 // last  post liked by a user
	SActionLikedPost     = "liked_post"
	IActionCommentedPost = 10 // last  post commented by a user
	SActionCommentedPost = "commented_post"
	IActionMentionedPost = 11 // last  post commented by a user
	SActionMentionedPost = "mentioned_post"
)

// const time to milliseconds
const (
	IOneMinute = 60000
	IOneHour   = 3600000
	IOneDay    = 86400000
	ITwoDays   = 172800000
	IOneWeek   = 604800000
	IOneMonth  = 2629746000
	SOneDay    = "86400000"
	STwoDays   = "172800000"
	SOneWeek   = "604800000"
)

// ErrorCode Table
const (
	EcSuccess                           = 1   //	Success
	EcNoExistObject                     = 2   //	No exist this object.
	EcExisObject                        = 3   //  Exist this object
	EcParam                             = 100 //	Invalid parameter
	EcParamMissingField                 = 101 //	Missing a few fields.
	EcParamUserID                       = 110 //	Invalid user id
	EcParamUserField                    = 111 //	Invalid user info field
	EcParamEmail                        = 113 //	Invalid email
	EcParamFieldList                    = 115 //	Invalid field list
	EcParamPhotoID                      = 121 //	Invalid photo id
	EcParamTitle                        = 142 //	Invalid title
	EcParamAccessToken                  = 190 //	Invalid OAuth 2.0 Access Token
	EcPermission                        = 200 //	Permissions error
	EcPermissionUser                    = 210 //	User not visible
	ECPermissionStatus                  = 220 //  Status not visible
	EcPermissionPhoto                   = 221 //	Photo not visible
	EcPermissionPost                    = 222 //  Post not visible
	EcPermissionGroup                   = 223 //  Group not visible
	EcPermissionMessage                 = 230 //	Permissions disallow message to user
	EcEdit                              = 300 //	Edit failure
	EcEditUserData                      = 310 //	User data edit failure
	EcUsersCreateInvalidEmail           = 370 //	The email address you provided is not a valid email address
	EcUsersCreateExistingEmail          = 371 //	The email address you provided belongs to an existing account
	EcUsersCreateBirthday               = 372 //	The birthday provided is not valid
	EcUsersCreatePassword               = 373 //	The password provided is too short or weak
	EcUsersRegisterInvalidCredential    = 374 //	The login credential you provided is invalid.
	EcUsersRegisterConfFailure          = 375 //	Failed to send confirmation message to the specified login credential.
	EcUsersRegisterExisting             = 376 //	The login credential you provided belongs to an existing account
	EcUsersRegisterDefaultError         = 377 //	Sorry, we were unable to process your registration.
	EcUsersRegisterPasswordBlank        = 378 //	Your password cannot be blank. Please try another.
	EcUsersRegisterPasswordInvalidChars = 379 //	Your password contains invalid characters. Please try another.
	EcUsersRegisterPasswordShort        = 380 //	Your password must be at least 6 characters long. Please try another.
	EcUsersRegisterPasswordWeak         = 381 //	Your password should be more secure. Please try another.
	EcUsersRegisterUsernameError        = 382 //	Please enter a valid username.
	EcUsersRegisterMissingInput         = 383 //	You must fill in all of the fields.
	EcUsersRegisterIncompleteBday       = 384 //	You must indicate your full birthday to register.
	EcUsersRegisterInvalidEmail         = 385 //	Please enter a valid email address.
	EcUsersRegisterEmailDisabled        = 386 //	The email address you entered has been disabled. Please contact disabled@facebook.com with any questions.
	EcUsersRegisterAddUserFailed        = 387 //	There was an error with your registration. Please try registering again.
	EcUsersRegisterNoGender             = 388 //	Please select either Male or Female.
	EcAuthEmail                         = 400 //	Invalid email address
	EcAuthLogin                         = 401 //	Invalid username or password
	EcAuthMissingToken                  = 404 //	Missing token.
	EcAuthInvalidToken                  = 405 //	Invalid token.
	EcAuthNoExistToken                  = 406 //	No exist token.
	EcAuthCheckToken                    = 407 //	Error in checking token.
	EcAuthGenerateToken                 = 408 //	Error in generate token.
	EcAuthNoExistUser                   = 409 //	No exist user.
	EcAuthNoExistFacebook               = 410 //	No exist account with this facebook.
	EcAuthInvalidFacebookToken          = 411 //	Error in checking token.
	EcAuthWrongPassword                 = 412 //	Error in login: Wrong password.
	EcAuthNoExistEmail                  = 413 //  No exist email
	EcAuthWrongRecoveryCode             = 414 //	Error in recover password: Wrong recovery code.
	EcMesgNoBody                        = 501 //	Missing message body
)

// ErrorCode  string Table
const (
	SEcSuccess                           = "Success"
	SEcNoExistObject                     = "No exist this object."
	SEcParam                             = "Invalid parameter"
	SEcParamMissingField                 = "Missing a few fields."
	SEcParamUserID                       = "Invalid user id"
	SEcParamUserField                    = "Invalid user info field"
	SEcParamEmail                        = "Invalid email"
	SEcParamFieldList                    = "Invalid field list"
	SEcParamPhotoID                      = "Invalid photo id"
	SEcParamTitle                        = "Invalid title" // ~needfix ~doing
	SEcParamAccessToken                  = 190             //	"Invalid OAuth 2.0 Access Token"
	SEcPermission                        = 200             //	"Permissions error"
	SEcPermissionUser                    = 210             //	"User not visible"
	SEcPermissionPhoto                   = 221             //	"Photo not visible"
	SEcPermissionMessage                 = 230             //	"Permissions disallow message to user"
	SEcEdit                              = 300             //	"Edit failure"
	SEcEditUserData                      = 310             //	"User data edit failure"
	SEcUsersCreateInvalidEmail           = 370             //	"The email address you provided is not a valid email address"
	SEcUsersCreateExistingEmail          = 371             //	"The email address you provided belongs to an existing account"
	SEcUsersCreateBirthday               = 372             //	"The birthday provided is not valid"
	SEcUsersCreatePassword               = 373             //	T"he password provided is too short or weak"
	SEcUsersRegisterInvalidCredential    = 374             //	The login credential you provided is invalid.
	SEcUsersRegisterConfFailure          = 375             //	Failed to send confirmation message to the specified login credential.
	SEcUsersRegisterExisting             = 376             //	The login credential you provided belongs to an existing account
	SEcUsersRegisterDefaultError         = 377             //	Sorry, we were unable to process your registration.
	SEcUsersRegisterPasswordBlank        = 378             //	Your password cannot be blank. Please try another.
	SEcUsersRegisterPasswordInvalidChars = 379             //	Your password contains invalid characters. Please try another.
	SEcUsersRegisterPasswordShort        = 380             //	Your password must be at least 6 characters long. Please try another.
	SEcUsersRegisterPasswordWeak         = 381             //	Your password should be more secure. Please try another.
	SEcUsersRegisterUsernameError        = 382             //	Please enter a valid username.
	SEcUsersRegisterMissingInput         = 383             //	You must fill in all of the fields.
	SEcUsersRegisterIncompleteBday       = 384             //	You must indicate your full birthday to register.
	SEcUsersRegisterInvalidEmail         = 385             //	Please enter a valid email address.
	SEcUsersRegisterEmailDisabled        = 386             //	The email address you entered has been disabled. Please contact disabled@facebook.com with any questions.
	SEcUsersRegisterAddUserFailed        = 387             //	There was an error with your registration. Please try registering again.
	SEcUsersRegisterNoGender             = 388             //	Please select either Male or Female.
	SEcAuthEmail                         = 400             //	Invalid email address
	SEcAuthLogin                         = 401             //	Invalid username or password
	SEcAuthMissingToken                  = 404             //	Missing token.
	SEcAuthInvalidToken                  = 405             //	Invalid token.
	SEcAuthNoExistToken                  = 406             //	No exist token.
	SEcAuthCheckToken                    = 407             //	Error in checking token.
	SEcAuthGenerateToken                 = 408             //	Error in generate token.
	SEcAuthNoExistUser                   = 409             //	No exist user.
	SEcAuthNoExistFacebook               = 410             //	No exist account with this facebook.
	SEcAuthInvalidFacebookToken          = 411             //	Error in checking token.
	SEcAuthWrongPassword                 = 412             //	Error in login: Wrong password.
	SEcAuthNoExistEmail                  = 413             //  No exist email
	SEcAuthWrongRecoveryCode             = 414             //	Error in recover password: Wrong recovery code.
	SEcMesgNoBody                        = 501             //	Missing message body
)

// TypePost const
const (
	IPost           = 0
	SPost           = "post"
	IPostStatus     = 1
	SPostStatus     = "status"
	IPostPhoto      = 2
	SPostPhoto      = "photo"
	IPostLink       = 3
	SPostLink       = "link"
	IPostGroup      = 4
	SPostGroup      = "post_group"
	IPostSharePost  = 5
	SPostSharePost  = "share_post"
	IPostSharePage  = 6
	SPostSharePage  = "share_page"
	IPostShareGroup = 7
	SPostShareGroup = "share_group"
)

// TypeNotification const
const (
	INotiPost    = 1
	SNotiPost    = "post"
	INotiFollow  = 2
	SNotiFollow  = "follow"
	INotiLike    = 3
	SNotiLike    = "like"
	INotiComment = 4
	SNotiComment = "comment"
	INotiStatus  = 5
	SNotiStatus  = "status"
	INotiPhoto   = 6
	SNotiPhoto   = "photo"
	INotiMention = 7
	SNotiMention = "mention"
	INotiAll     = 8
	SNotiAll     = "all"
)

// Role user
const (
	IUserRole       = 1
	IStudentRole    = 2
	ITeacherRole    = 3
	ISupervisorRole = 4
	IAdminRole      = 5
)

//FCMToken struct
const (
	FCMToken = "AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V"
)

// WS from TLU
const (
	SURLGetSemesterListByYear             = "https://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?DanhSachHocKy&NamHoc="
	SURLGetSubjectListBySemesterCode      = "https://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?DanhSachHocPhan&HocKy="
	SURLGetTeacherListBySemesterCode      = "https://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?DanhSachGiangVien&HocKy="
	SURLGetClassListBySemesterCode        = "https://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?DanhSachLopHocPhan&HocKy="
	SURLGetStudentListByClassCode         = "https://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?DanhSachSinhVienLopHocPhan&MaLop="
	SURLGetExamScheduleListBySemesterCode = "http://elearning.thanglong.edu.vn/tlu-custom/ws/tluscn.ws.php?LichThi&HocKy="
)
