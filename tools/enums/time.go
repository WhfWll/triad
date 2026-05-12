package enums

const (
	// TimeLayout 完整时间返回
	TimeLayout = "2006-01-02 15:04:05"
	// TimeToMinLayout 只展示到分钟位置
	TimeToMinLayout = "2006-01-02 15:04:00"
	// TimeToHourLayout 只展示到小时位置
	TimeToHourLayout = "2006-01-02 15:00:00"
	// TimeToDayLayout 只展示天的位置
	TimeToDayLayout = "2006-01-02 00:00:00"
	// TimeYMDLayout 年月日时间
	TimeYMDLayout = "2006.01.02"

	//格式：年-月
	TimeYMonthLayout = "2006-01"

	//格式：年-月-日
	TimeYMDBarLayout = "2006-01-02"

	//格式：年-月-日 时
	TimeYMDHourLayout = "2006-01-02 15"

	//格式：年-月-日 时:分
	TimeYMDHMinLayout = "2006-01-02 15:04"

	// ResTimeMinLayout 分钟时间 - 接口返回需要的时间格式
	ResTimeMinLayout = "15:04"
	// ResTimeHHourLayout 半小时时间
	ResTimeHHourLayout = "01-02 15:04"

	// ResTimeHHour0Layout 半小时时间 0分
	ResTimeHHour0Layout = "01-02 15:00"

	// ResTimeHHour30Layout 半小时时间 30分
	ResTimeHHour30Layout = "01-02 15:00"

	// ResTimeDayLayout 天时间
	ResTimeDayLayout = "2006-01-02"
	// ResTimeSecLayout 10s时间 - 接口返回需要的时间格式
	ResTimeSecLayout = "04"
)
