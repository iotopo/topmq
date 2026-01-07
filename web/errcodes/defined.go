package errcodes

import (
	"net/http"

	"github.com/xuri/excelize/v2"
)

var (
	// SessionIPChanged 用户 IP 与会话不一致
	SessionIPChanged = NewHttpError(http.StatusUnauthorized, "session ip changed")
	// SessionExpired 由于用户被禁用、用户密码修改等原因，强制下线
	SessionExpired = NewHttpError(http.StatusUnauthorized, "session expired")
	// InvalidCsrfToken csrf token 无效
	InvalidCsrfToken = NewHttpError(http.StatusUnauthorized, "invalid csrf token")
	// 登录：验证码错误
	BadCaptcha = NewSimpleFailure("bad-captcha")
	// 登录：用户名或密码错误
	BadUserOrPassword = NewSimpleFailure("bad-user-or-password")
	// 登录：用户已被锁定
	UserLocked = NewSimpleFailure("user-locked")

	// NameDuplicated  名称重复
	NameDuplicated = NewSimpleFailure("err_name_duplicated")
	// NotFound        记录不存在或已被删除
	NotFound = NewSimpleFailure("err_not_found")

	// SQLite 专用，因为 SQlite 在外键错误时不会提示外键名，因为无法得知是哪一个表的外键
	ForeignKeyError      = NewSimpleFailure("foreign-key")
	ForeignRestrictError = NewSimpleFailure("foreign-restrict")

	// 导入：文件过大
	ImportFileTooLarge = NewSimpleFailure("import-file-too-large")
	// 导入：无法打开文件（文件损坏或格式错误）
	ImportFileBroken = NewSimpleFailure("import-file-broken")
)

// 导入：Excel工作表不存在
type ImportSheetNotFound struct {
	Sheet string
}

func (e ImportSheetNotFound) Build() ErrorResponse {
	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_sheet_not_found",
		ErrorMsg:   "sheet does not exist",
		ErrorArgs:  map[string]any{"sheetName": e.Sheet},
	}
}

// 导入 Excel 通用
// 表 ${sheet} 中的单元格 ${cell} 必填
// 如果 ${refCell} 存在，表示此单元格并非始终必填，是由于 ${refCell} 填写的值所以必填
type ImportCellRequired struct {
	Sheet  string
	Row    int
	Col    int
	RefCol int

	ColName string // col == 0 时替代
}

func (e ImportCellRequired) Build() ErrorResponse {
	args := map[string]any{
		"sheet": e.Sheet,
	}
	if e.Col > 0 {
		args["cell"], _ = excelize.CoordinatesToCellName(e.Col, e.Row)
	} else {
		args["cell"] = e.ColName
	}
	if e.RefCol > 0 {
		args["refCell"], _ = excelize.CoordinatesToCellName(e.RefCol, e.Row)
	}

	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_cell_required",
		ErrorMsg:   "cell value is required",
		ErrorArgs:  args,
	}
}

// 导入 Excel 通用
// 表 ${sheet} 中的单元格 ${cell} 填写错误
// 如果 ${refCell} 存在，表示此单元格的值并非始终填写错误，是由于 ${refCell} 填写的值所以填写错误
type ImportCellIllegal struct {
	Sheet  string
	Row    int
	Col    int
	RefCol int
}

func (e ImportCellIllegal) Build() ErrorResponse {
	args := map[string]any{
		"sheet": e.Sheet,
	}
	args["cell"], _ = excelize.CoordinatesToCellName(e.Col, e.Row)
	if e.RefCol > 0 {
		args["refCell"], _ = excelize.CoordinatesToCellName(e.RefCol, e.Row)
	}

	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_cell_illegal",
		ErrorMsg:   "cell value is illegal",
		ErrorArgs:  args,
	}
}

// 导入 Excel 通用
// 表 ${sheet} 中的单元格 ${cell} 引用错误（引用了不存在的设备模型等）
type ImportCellRefError struct {
	Sheet string
	Row   int
	Col   int
	Value string // (optional) 单元格中有多个值时，返回出现引用错误的值
}

func (e ImportCellRefError) Build() ErrorResponse {
	args := map[string]any{
		"sheet": e.Sheet,
	}
	args["cell"], _ = excelize.CoordinatesToCellName(e.Col, e.Row)
	if e.Value != "" {
		args["value"] = e.Value
	}

	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_cell_ref_error",
		ErrorMsg:   "cell value reference error",
		ErrorArgs:  args,
	}
}

// 导入 Excel 通用
// 表 ${sheet} 中的单元格 ${cell} 值重复
type ImportCellDuplicated struct {
	Sheet string
	Row   int
	Col   int
}

func (e ImportCellDuplicated) Build() ErrorResponse {
	args := map[string]any{
		"sheet": e.Sheet,
	}
	args["cell"], _ = excelize.CoordinatesToCellName(e.Col, e.Row)

	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_cell_duplicated",
		ErrorMsg:   "cell value duplicated",
		ErrorArgs:  args,
	}
}

// 导入 Excel 通用
// 表 ${sheet} 中，由于单元格 ${refCell} 填写的值，单元格 ${cell} 必须为空
type ImportCellForbidden struct {
	Sheet  string
	Row    int
	Col    int
	RefCol int
}

func (e ImportCellForbidden) Build() ErrorResponse {
	args := map[string]any{
		"sheet": e.Sheet,
	}
	args["cell"], _ = excelize.CoordinatesToCellName(e.Col, e.Row)
	args["refCell"], _ = excelize.CoordinatesToCellName(e.RefCol, e.Row)

	return ErrorResponse{
		StatusCode: http.StatusOK,
		ErrorCode:  "err_import_cell_forbidden",
		ErrorMsg:   "cell value is forbidden",
		ErrorArgs:  args,
	}
}
