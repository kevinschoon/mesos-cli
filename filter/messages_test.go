package filter

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

const (
	stateSize = 25000
	totalMsgs = stateSize * 4
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type chooser struct {
	state map[int]bool
}

func (c *chooser) int() int {
	choice := rand.Intn(stateSize)
	if _, ok := c.state[choice]; ok {
		return c.int()
	}
	c.state[choice] = true
	return choice
}

func mockState(kts []proto.Message, t *testing.T) Messages {
	messages := Messages{}
	for i := 0; i < stateSize; i++ {
		messages = append(messages, &mesos.Task{
			TaskID: mesos.TaskID{Value: RandString(25)},
		})
		messages = append(messages, &mesos.ExecutorInfo{
			ExecutorID: mesos.ExecutorID{Value: RandString(25)},
		})
		messages = append(messages, &mesos.FrameworkInfo{
			ID: &mesos.FrameworkID{Value: RandString(25)},
		})
		messages = append(messages, &mesos.AgentInfo{
			ID: &mesos.AgentID{Value: RandString(25)},
		})
	}
	// Randomly distribute the supplied known types in our state for testing
	choices := chooser{map[int]bool{}}
	for _, kt := range kts {
		messages[choices.int()] = kt
	}
	if t != nil {
		fmt.Println(len(messages))
		assert.Len(t, messages, totalMsgs)
	}
	return messages
}

func TestFindAny(t *testing.T) {

	messages := mockState(Messages{
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task"}},
		&mesos.ExecutorInfo{ExecutorID: mesos.ExecutorID{Value: "test-executor"}},
		&mesos.FrameworkInfo{ID: &mesos.FrameworkID{Value: "test-framework"}},
		&mesos.AgentInfo{ID: &mesos.AgentID{Value: "test-agent"}},
	}, t)

	// TaskID
	_, err := messages.FindAny(TaskIDFilter("not-real", false))
	assert.Error(t, err, "Not Found")
	task, err := AsTask(messages.FindAny(TaskIDFilter("test-task", false)))
	assert.NoError(t, err)
	assert.Equal(t, task.TaskID.Value, "test-task")
	task, err = AsTask(messages.FindAny(TaskIDFilter("test-", true)))
	assert.NoError(t, err)
	assert.Equal(t, task.TaskID.Value, "test-task")

	// ExecutorID
	_, err = messages.FindAny(ExecutorIDFilter("not-real", false))
	assert.Error(t, err, "Not Found")
	executor, err := AsExecutor(messages.FindAny(ExecutorIDFilter("test-executor", false)))
	assert.NoError(t, err)
	assert.Equal(t, executor.ExecutorID.Value, "test-executor")
	executor, err = AsExecutor(messages.FindAny(ExecutorIDFilter("test-", true)))
	assert.NoError(t, err)
	assert.Equal(t, executor.ExecutorID.Value, "test-executor")

	// FrameworkID
	_, err = messages.FindAny(FrameworkIDFilter("not-real", false))
	assert.Error(t, err, "Not Found")
	framework, err := AsFramework(messages.FindAny(FrameworkIDFilter("test-framework", false)))
	assert.NoError(t, err)
	assert.Equal(t, framework.ID.Value, "test-framework")
	framework, err = AsFramework(messages.FindAny(FrameworkIDFilter("test-", true)))
	assert.NoError(t, err)
	assert.Equal(t, framework.ID.Value, "test-framework")

	// AgentID
	_, err = messages.FindAny(AgentIDFilter("not-real", false))
	assert.Error(t, err, "Not Found")
	agent, err := AsAgent(messages.FindAny(AgentIDFilter("test-agent", false)))
	assert.NoError(t, err)
	assert.Equal(t, agent.ID.Value, "test-agent")
	agent, err = AsAgent(messages.FindAny(AgentIDFilter("test-", true)))
	assert.NoError(t, err)
	assert.Equal(t, agent.ID.Value, "test-agent")
}

func TestFindOne(t *testing.T) {
	messages := mockState(Messages{
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-1"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-2"}},
	}, nil)
	_, err := messages.FindOne(TaskIDFilter("test-task-", true))
	assert.Error(t, err, "Too many results")
	_, err = messages.FindOne(TaskIDFilter("not-real", false))
	assert.Error(t, err, "Not found")
	task, err := AsTask(messages.FindOne(TaskIDFilter("test-task-1", false)))
	assert.NoError(t, err)
	assert.Equal(t, task.TaskID.Value, "test-task-1")
}

func TestFindMany(t *testing.T) {
	messages := mockState(Messages{
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-0"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-1"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-2"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-3"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-4"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-5"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-6"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-7"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-8"}},
		&mesos.Task{TaskID: mesos.TaskID{Value: "test-task-9"}},
	}, nil)
	tasks := AsTasks(messages.FindMany(TaskIDFilter("test-task-", true)))
	assert.Len(t, tasks, 10)
}

func BenchmarkMessagesWorst(b *testing.B) {
	messages := mockState(Messages{}, nil)
	// Benchmark worst case scenario
	last := messages[len(messages)-1]
	var filter Filter
	switch t := last.(type) {
	case *mesos.Task:
		filter = TaskIDFilter(t.TaskID.Value, false)
	case *mesos.ExecutorInfo:
		filter = ExecutorIDFilter(t.ExecutorID.Value, false)
	case *mesos.FrameworkInfo:
		filter = FrameworkIDFilter(t.ID.Value, false)
	case *mesos.AgentInfo:
		filter = AgentIDFilter(t.ID.Value, false)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		messages.FindAny(filter)
	}
}

func BenchmarkMessagesBest(b *testing.B) {
	messages := mockState([]proto.Message{}, nil)
	// Benchmark worst case scenario
	last := messages[0]
	var filter Filter
	switch t := last.(type) {
	case *mesos.Task:
		filter = TaskIDFilter(t.TaskID.Value, false)
	case *mesos.ExecutorInfo:
		filter = ExecutorIDFilter(t.ExecutorID.Value, false)
	case *mesos.FrameworkInfo:
		filter = FrameworkIDFilter(t.ID.Value, false)
	case *mesos.AgentInfo:
		filter = AgentIDFilter(t.ID.Value, false)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		messages.FindAny(filter)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
