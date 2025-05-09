package job

import (
	"fmt"
	"sen-global-api/pkg/monitor"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

var instantiated *TimeMachine = nil
var once sync.Once

func New() *TimeMachine {
	once.Do(func() {
		instantiated = new(TimeMachine)
		instantiated.formExecutors = make([]IntervalTaskExecutor, 0)
		instantiated.form2Executors = make([]IntervalTaskExecutor, 0)
		instantiated.form3Executors = make([]IntervalTaskExecutor, 0)
		instantiated.form4Executors = make([]IntervalTaskExecutor, 0)
		instantiated.urlExecutors = make([]IntervalTaskExecutor, 0)
		instantiated.deviceSyncExecutors = make([]IntervalTaskExecutor, 0)
		instantiated.todoExecutors = make([]IntervalTaskExecutor, 0)
		instantiated.googleQPIRequestMonitor = make([]IntervalTaskExecutor, 0)
		instantiated.submissionSyncExecutors = make([]IntervalTaskExecutor, 0)
		instantiated.formCron = gocron.NewScheduler(time.UTC)
		instantiated.form2Cron = gocron.NewScheduler(time.UTC)
		instantiated.form3Cron = gocron.NewScheduler(time.UTC)
		instantiated.form4Cron = gocron.NewScheduler(time.UTC)
		instantiated.urlCron = gocron.NewScheduler(time.UTC)
		instantiated.deviceSyncCron = gocron.NewScheduler(time.UTC)
		instantiated.todoCron = gocron.NewScheduler(time.UTC)
		instantiated.googleQPIRequestMonitorCron = gocron.NewScheduler(time.UTC)
		instantiated.submissionSyncCron = gocron.NewScheduler(time.UTC)
	})
	return instantiated
}

type TimeMachine struct {
	formExecutors               []IntervalTaskExecutor
	form2Executors              []IntervalTaskExecutor
	form3Executors              []IntervalTaskExecutor
	form4Executors              []IntervalTaskExecutor
	urlExecutors                []IntervalTaskExecutor
	deviceSyncExecutors         []IntervalTaskExecutor
	todoExecutors               []IntervalTaskExecutor
	googleQPIRequestMonitor     []IntervalTaskExecutor
	submissionSyncExecutors     []IntervalTaskExecutor
	formCron                    *gocron.Scheduler
	form2Cron                   *gocron.Scheduler
	form3Cron                   *gocron.Scheduler
	form4Cron                   *gocron.Scheduler
	urlCron                     *gocron.Scheduler
	deviceSyncCron              *gocron.Scheduler
	todoCron                    *gocron.Scheduler
	googleQPIRequestMonitorCron *gocron.Scheduler
	submissionSyncCron          *gocron.Scheduler
}

type IntervalTaskExecutor interface {
	ExecuteSyncForms()
	ExecuteSyncForms2()
	ExecuteSyncForms3()
	ExecuteSyncForms4()
	ExecuteSyncUrls()
	ExecuteSyncTodos()
	ExecuteGoogleAPIRequestMonitor()
}

func (receiver *TimeMachine) Start(formInterval uint64, urlInterval uint64, todoInterval uint64, formInterval2 uint64, formInterval3 uint64, formInterval4 uint64) {
	receiver.ScheduleSyncForms(formInterval)
	receiver.ScheduleSyncForms2(formInterval2)
	receiver.ScheduleSyncForms3(formInterval3)
	receiver.ScheduleSyncForms4(formInterval4)
	receiver.ScheduleSyncUrls(urlInterval)
	receiver.ScheduleSyncToDos(todoInterval)
	receiver.ScheduleGoogleAPIRequestMonitor()

	monitor.SendMessageViaTelegram("Time machine started with ",
		fmt.Sprint("formInterval: ", formInterval),
		fmt.Sprint("formInterval2: ", formInterval2),
		fmt.Sprint("formInterval3: ", formInterval3),
		fmt.Sprint("formInterval4: ", formInterval4),
		fmt.Sprint("urlInterval: ", urlInterval),
		fmt.Sprint("todoInterval: ", todoInterval),
	)
}

func (receiver *TimeMachine) Stop() {
	receiver.formCron.Clear()
	receiver.form2Cron.Clear()
	receiver.form3Cron.Clear()
	receiver.form4Cron.Clear()
	receiver.urlCron.Clear()
	receiver.deviceSyncCron.Clear()
	receiver.todoCron.Clear()
	receiver.googleQPIRequestMonitorCron.Clear()
	receiver.submissionSyncCron.Clear()

	monitor.SendMessageViaTelegram("Time machine has been stopped")
}

func (receiver *TimeMachine) SubscribeFormsExec(exec IntervalTaskExecutor) {
	receiver.formExecutors = append(receiver.formExecutors, exec)
	log.Debug("Subscribe form executor", receiver.formExecutors)
}

func (receiver *TimeMachine) SubscribeForms2Exec(exec IntervalTaskExecutor) {
	receiver.form2Executors = append(receiver.form2Executors, exec)
	log.Debug("Subscribe form2 executor", receiver.form2Executors)
}

func (receiver *TimeMachine) SubscribeForms3Exec(exec IntervalTaskExecutor) {
	receiver.form3Executors = append(receiver.form3Executors, exec)
	log.Debug("Subscribe form3 executor", receiver.form3Executors)
}

func (receiver *TimeMachine) SubscribeForms4Exec(exec IntervalTaskExecutor) {
	receiver.form4Executors = append(receiver.form4Executors, exec)
	log.Debug("Subscribe form4 executor", receiver.form4Executors)
}

func (receiver *TimeMachine) SubscribeUrlsExec(exec IntervalTaskExecutor) {
	receiver.urlExecutors = append(receiver.urlExecutors, exec)
	log.Debug("Subscribe url executor", receiver.urlExecutors)
}

func (receiver *TimeMachine) SubscribeSyncDevicesExec(exec IntervalTaskExecutor) {
	receiver.deviceSyncExecutors = append(receiver.deviceSyncExecutors, exec)
	log.Debug("Subscribe device sync exec", receiver.deviceSyncExecutors)
}

func (receiver *TimeMachine) SubscribeSyncToDosExec(exec IntervalTaskExecutor) {
	receiver.todoExecutors = append(receiver.todoExecutors, exec)
	log.Debug("Subscribe todo sync exec", receiver.todoExecutors)
}

func (receiver *TimeMachine) SubscribeGoogleAPIRequestMonitorExec(exec IntervalTaskExecutor) {
	receiver.googleQPIRequestMonitor = append(receiver.googleQPIRequestMonitor, exec)
	log.Debug("Subscribe google api request monitor exec", receiver.googleQPIRequestMonitor)
}

func (receiver *TimeMachine) ScheduleSyncForms(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.formCron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.formCron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync forms")
		for _, executor := range receiver.formExecutors {
			executor.ExecuteSyncForms()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync formCron every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.formCron.StartAsync()
}

func (receiver *TimeMachine) ScheduleSyncForms2(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.form2Cron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.form2Cron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync forms")
		for _, executor := range receiver.form2Executors {
			executor.ExecuteSyncForms2()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync form2Cron every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.form2Cron.StartAsync()
}

func (receiver *TimeMachine) ScheduleSyncForms3(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.form3Cron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.form3Cron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync forms")
		for _, executor := range receiver.form3Executors {
			executor.ExecuteSyncForms3()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync form3Cron every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.form3Cron.StartAsync()
}

func (receiver *TimeMachine) ScheduleSyncForms4(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.form4Cron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.form4Cron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync forms")
		for _, executor := range receiver.form4Executors {
			executor.ExecuteSyncForms4()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync form4Cron every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.form4Cron.StartAsync()
}

func (receiver *TimeMachine) ScheduleSyncUrls(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.urlCron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.urlCron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync urls")
		for _, executor := range receiver.urlExecutors {
			executor.ExecuteSyncUrls()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync url every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.urlCron.StartAsync()
}

func (receiver *TimeMachine) ScheduleSyncToDos(interval uint64) {
	if interval == 0 {
		return
	}
	receiver.todoCron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(interval) * time.Minute)
	task, err := receiver.todoCron.Every(int(interval)).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Sync todos")
		for _, executor := range receiver.todoExecutors {
			executor.ExecuteSyncTodos()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule sync todos every ", interval, " minutes [ERROR]? ", task.Error())
	}
	receiver.todoCron.StartAsync()
}

func (receiver *TimeMachine) ScheduleGoogleAPIRequestMonitor() {
	receiver.googleQPIRequestMonitorCron.Clear()

	now := time.Now()
	startAt := now.Add(time.Duration(1) * time.Minute)
	task, err := receiver.googleQPIRequestMonitorCron.Every(1).Minutes().StartAt(startAt).Do(func() {
		log.Debug("Report Google API Request Monitor")
		for _, executor := range receiver.googleQPIRequestMonitor {
			executor.ExecuteGoogleAPIRequestMonitor()
		}
	})
	if err != nil {
		log.Error(err)
		panic(err)
	} else if task.Error() != nil {
		log.Error(task.Error())
		panic(task.Error())
	} else if task != nil && task.Error() == nil {
		log.Info("Schedule monitor google api request [ERROR]? ", task.Error())
	}
	receiver.googleQPIRequestMonitorCron.StartAsync()
}
