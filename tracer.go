package bingo

/*****************************************************/
/* 实现自定义的context接口                             */
/* 作用：使context和cron能够使用backend和db、redis包等  */
/*****************************************************/
type Tracer interface {
	Logs(args ...string)
}

