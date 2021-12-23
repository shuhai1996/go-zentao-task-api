package zentao

// service.go:
import (
	"fmt"
	"go-zentao-task/model/zentao"
	"strconv"
	"time"
)

func InitializeService() *Service {
	task := zentao.NewTask()
	user := zentao.NewUser()
	estimate := zentao.NewTaskEstimate()
	projectProduct := zentao.NewProjectProduct()
	action := zentao.NewAction()
	service := NewService(task, user, estimate, projectProduct, action)
	return service
}

func NewService(
	task *zentao.Task, user *zentao.User, estimate *zentao.TaskEstimate, projectProduct *zentao.ProjectProduct, action *zentao.Action) *Service {
	return &Service{
		Task:           task,
		User:           user,
		Estimate:       estimate,
		ProjectProduct: projectProduct,
		Action:         action,
	}
}

type Service struct {
	Task           *zentao.Task
	User           *zentao.User
	Estimate       *zentao.TaskEstimate
	ProjectProduct *zentao.ProjectProduct
	Action         *zentao.Action
}

func (service *Service) TaskView(id int) string {
	//获取任务
	t, err := service.Task.FindOne(id)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if t.ID == 0 {
		fmt.Println("任务不存在")
		return ""
	}
	return t.Name
}

// GetAllTaskNotDone 获取未完成任务
func (service *Service) GetAllTaskNotDone() []zentao.Task {
	tasks, err := service.Task.FindAll(service.User.Account, "", 1, "asc")
	if err != nil {
		fmt.Println("获取任务异常")
		return nil
	}
	return tasks
}

// GetEstimateToday 获取今日工时
func (service *Service) GetEstimateToday() (float64, error) {
	estimate, err := service.Estimate.GetToday(service.User.Account)
	if err != nil {
		fmt.Println("获取工时异常")
		return 0, nil
	}
	var consumed float64
	for _, v := range estimate {
		consumed += v.Consumed
	}
	fmt.Println("今日工时填写", fmt.Sprintf("%.2f", consumed))
	return consumed, err
}

func (service *Service) UpdateTask(task int, estimate float64, action string) float64 {
	taskInfo, err := service.Task.FindOne(task)
	var name = ""
	if err != nil {
		fmt.Println("获取任务异常")
		return 0
	}
	productInfo, _ := service.ProjectProduct.FindOneByProject(taskInfo.Project)

	if taskInfo.Status == zentao.StatusWait { // 等待状态改为开始
		action = zentao.ActionStart
		taskInfo.Status = zentao.StatusDo
	}

	if estimate == taskInfo.Left || estimate > taskInfo.Left { //消耗工时不能大于剩余工时
		estimate = taskInfo.Left
		action = zentao.StatusFinished
		taskInfo.Status = zentao.StatusFinished
	}
	if action == zentao.StatusFinished {
		if taskInfo.FromBug != 0 {
			fmt.Println("任务从bug创建，不能直接完成")
			return 0
		}
		estimate = taskInfo.Left
		taskInfo.FinishedDate = time.Now()
		name = taskInfo.AssignedTo
	}

	if estimate <= 0 {
		fmt.Println("耗时不能为0")
		return 0
	}

	//创建操作记录
	service.Action.Create(task, "task", ","+strconv.Itoa(productInfo.Product)+",", taskInfo.Project, taskInfo.Execution, taskInfo.AssignedTo, action, strconv.FormatFloat(estimate, 'f', -1, 64))
	//创建工时填写记录
	service.Estimate.Create(task, taskInfo.Left, estimate, taskInfo.AssignedTo, "")
	//更新任务
	service.Task.UpdateOne(task, estimate+taskInfo.Consumed, taskInfo.Left-estimate, name, taskInfo.FinishedDate, taskInfo.Status)
	return estimate
}

func (service *Service) GetOptimumTasks() []int {
	var tasks []zentao.Task
	var res []int
	return service.OptimumTask(tasks, res, 0)
}

// OptimumTask 竞选出最优任务， 用于更新
func (service *Service) OptimumTask(tasks []zentao.Task, result []int, round int) []int {
	var tmpTasks []zentao.Task
	if len(tasks) == 0 {
		tasks = service.GetAllTaskNotDone()
		result = []int{}
		round = 0
	}
	if round > 3 {
		return result
	}
	for _, task := range tasks {
		switch task.Type {
		case zentao.TypeDiscuss: //优先讨论类型的任务
			result = append(result, task.ID)
		case zentao.TypeDev:
			if round == 1 && task.Left <= 8 && task.Status == zentao.StatusDo {
				result = append(result, task.ID)
			} else if round == 2 && task.Left > 8 && task.Status == zentao.StatusDo { // 第三轮竞选，剩余时间大于8天并且在doing状态的
				result = append(result, task.ID)
			} else if round == 3 && task.Status == zentao.StatusWait && task.Left > 0 { // 第四轮竞选，等待状态的case,并且剩余时间不能为0
				result = append(result, task.ID)
			} else {
				tmpTasks = append(tmpTasks, task)
			}
		}
	}
	if len(tmpTasks) > 0 {
		result = service.OptimumTask(tmpTasks, result, round+1)
	}
	return result
}

// ConsumeRecord 记录工时
func (service *Service) ConsumeRecord(estimate float64) (float64, []int) {
	var current float64
	var ids []int
	if estimate <= 0 {
		fmt.Println("今日工时已完成，不需要填写")
		return 0, nil
	}
	tasks := service.GetOptimumTasks()
OuterLoop: // 循环标签
	for _, task := range tasks {
		if estimate-current > 0 {
			count := service.UpdateTask(task, estimate-current, zentao.ActionRecorde)
			current += count
			if count > 0 {
				ids = append(ids, task)
			}
			if current >= 8 {
				break OuterLoop //跳出循环
			}
		}
	}
	current, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", current), 64)
	return current, ids
}

func (service *Service) UserLogin() {
	action, _ := service.Action.FindLastLogin(service.User.Account)
	// 当前时间往前推三个小时
	add, _ := time.ParseDuration("-3h")
	last := time.Now().Add(add)
	if last.Sub(action.Date).Seconds() > 0 { //记录时间在三小时前
		//创建操作记录
		service.Action.Create(service.User.ID, "user", ","+strconv.Itoa(0)+",", 0, 0, service.User.Account, zentao.ActionLogin, "")
	}
}
