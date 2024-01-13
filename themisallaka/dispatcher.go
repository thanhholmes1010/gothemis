package themisallaka

type DispatcherManager struct {
	*pool
	// underlying is thread pool
	// main job is send letter to mailbox
	// if not exist job, then sleep, release cycle cpu
}

type ExecutorManager struct {
	*pool
	// underlying is threadpool
	// main job is execute mail box
	// if not exist job, then sleep, release cycle cpu
}

// son gui cho thanh.process, "nay di nhau khong"

// t1: nhung son se gui cho buudien(=manager) truoc

//     t2: dominic se gui cho thanh.process "hello"
//     t2: nhung dominic se gui cho buudien(=manager) truoc

// neu buu dien chi co 1 cong nhan thi sao ?
// thi t1 di gui den thanh.process (= hop thu truoc cua nha thanh`), hop thu nay se la queue order

// o thoi diem t3: tui lay t1 ra  (nguoi doc 1)
//     o thoi diem t4: tui lay t2 ra (nguoi doc 2)

// thi neu nhu buc thu t1 nay co thoi gian thuc hien lau hon, thi no khong bi bo doi' buc thu t2
