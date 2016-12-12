package base

import (
    //linuxproc "github.com/c9s/goprocinfo/linux"
)
import (
    "time"
    "sync"
    "fmt"
    "strconv"
	"gatoor/orca/rewriteTrainer/needs"

)

const (
    APP_HTTP = "http"
    APP_WORKER = "worker"

    STATUS_INIT = "init"
    STATUS_RUNNING = "running"
    STATUS_DEPLOYING = "deploying"
    STATUS_DEAD = "dead"
    STATUS_UNKNOWN = "unknown"

    FILE_COMMAND = "FILE_COMMAND"
    EXEC_COMMAND = "EXEC_COMMAND"
)

type HostId string
type Version uint64
type IpAddr string
type HabitatName string
type HabitatStatus string
type Status string
type DeploymentCount int
type LoadBalancerName string
type NetworkName string

type Command struct {
    Path string
    Args string
}

type OsCommandType string

type OsCommand struct {
    Type OsCommandType
    Command Command
}

type Usage uint64

type HostStats struct {
    MemoryUsage Usage
    CpuUsage Usage
    NetworkUsage Usage
}

type AppStats struct {
    MemoryUsage Usage
    CpuUsage Usage
    NetworkUsage Usage
    ResponsePerformance uint64
}

type HostInfo struct {
    HostId HostId
    IpAddr IpAddr
    OsInfo OsInfo
    Apps []AppInfo
}

type OsInfo struct {
    Os string
    Arch string
}

type HabitatInfo struct {
    Version Version
    Name HabitatName
    Status Status
}

type AppInfo struct {
    Type AppType
    Name AppName
    Version Version
    Status Status
    Id AppId
}

type StatsWrapper struct {
    Host HostStats
    Apps map[AppName]AppStats
}

type TrainerPushWrapper struct {
    HostInfo HostInfo
    Stats MetricsWrapper
}

type AppMetrics map[AppName]map[Version]map[string]AppStats
type AppMetricsJson map[AppName]map[string]map[string]AppStats

type MetricsWrapper struct {
    HostMetrics map[string]HostStats
    AppMetrics AppMetricsJson
}

var metricsMutex = &sync.Mutex{}

func (m *MetricsWrapper) Wipe() {
    metricsMutex.Lock()
    defer metricsMutex.Unlock()
    (*m).HostMetrics = make(map[string]HostStats)
    (*m).AppMetrics = make(map[AppName]map[string]map[string]AppStats)
}

func (m MetricsWrapper) Get() MetricsWrapper{
    metricsMutex.Lock()
    defer metricsMutex.Unlock()
    metrics := m
    return metrics
}

func (m MetricsWrapper) AddHostMetrics(hostMetrics HostStats) {
    metricsMutex.Lock()
    defer metricsMutex.Unlock()
    m.HostMetrics[time.Now().UTC().Format(time.RFC3339Nano)] = hostMetrics
}

func (m MetricsWrapper) AddAppMetrics(appName AppName, version Version, appMetrics AppStats) {
    metricsMutex.Lock()
    defer metricsMutex.Unlock()
    versionStr := strconv.Itoa(int(version))
    if _, exists := m.AppMetrics[appName]; !exists {
        m.AppMetrics[appName] = make(map[string]map[string]AppStats)
    }
    if _, exists := m.AppMetrics[appName][versionStr]; !exists {
        m.AppMetrics[appName][versionStr] = make(map[string]AppStats)
    }
    m.AppMetrics[appName][versionStr][time.Now().UTC().Format(time.RFC3339Nano)] = appMetrics
}


type PushWrapper struct {
    HostInfo HostInfo
    Stats StatsWrapper
}

type PushConfiguration struct {
    //TargetHostId HostId
    OrcaVersion string
    DeploymentCount DeploymentCount
    AppConfiguration AppConfiguration
    //HabitatConfiguration HabitatConfiguration
}

type HabitatConfiguration struct {
    Name HabitatName
    Version Version
    InstallCommands []OsCommand
}

type AppConfiguration struct {
    Name                  AppName
    Type                  AppType
    Version               Version
    TargetDeploymentCount DeploymentCount
    MinDeploymentCount    DeploymentCount
    DockerConfig          DockerConfig
    RawConfig             RawConfig
    LoadBalancer LoadBalancerName
    Network NetworkName
    Needs needs.AppNeeds
}

type ProviderType string
type InstanceType string
type Cost float32
type InstanceCount int
type SafeInstance bool

type Resource int
type CpuResource Resource
type MemoryResource Resource
type NetworkResource Resource

type InstanceResources struct {
    UsedCpuResource CpuResource
    UsedMemoryResource MemoryResource
    UsedNetworkResource NetworkResource
    TotalCpuResource CpuResource
    TotalMemoryResource MemoryResource
    TotalNetworkResource NetworkResource
}

type AWSConfiguration struct {
    Key string
    Secret string
    Region string
    AMI string
    SecurityGroupId string

    // Transient Fields that are populated by the system when the aws provider is inited
    InstanceTypes []InstanceType
    InstanceCost map[InstanceType]Cost
    InstanceResources map[InstanceType]InstanceResources
    InstanceSafety map[InstanceType]SafeInstance
    SuitableInstanceSafetyFactor float32
}


type ProviderConfiguration struct {
    Type ProviderType
    SSHKey string
    SSHUser string
    MinInstances InstanceCount
    MaxInstances InstanceCount
    AWSConfiguration AWSConfiguration
}

type RawConfig struct {
    InstallCommands []OsCommand
    RunCommand OsCommand
    QueryStateCommand OsCommand
    RemoveCommand OsCommand
    StopCommand OsCommand
}

type DockerConfig struct {
    Tag string
    Repository string
    Reference string
}

type TrainerPolicies struct {
    TRY_TO_REMOVE_HOSTS bool
}

type TrainerConfigurationState struct {
    Port int
    Policies TrainerPolicies
    Ip IpAddr
}

var appsMetricsMutex = &sync.Mutex{}

func (a *AppMetrics) Add(name AppName, version Version, time string, stats AppStats) {
    appsMetricsMutex.Lock()
    defer appsMetricsMutex.Unlock()
    if _, exists := (*a)[name]; !exists {
        (*a)[name] = make(map[Version]map[string]AppStats)
    }
    if _, exists := (*a)[name][version]; !exists {
        (*a)[name][version] = make(map[string]AppStats)
    }
    (*a)[name][version][time] = stats
}

func (a *AppMetrics) All() map[AppName]map[Version]map[string]AppStats {
    appsMetricsMutex.Lock()
    defer appsMetricsMutex.Unlock()
    res := (*a)
    return res
}

func (a *AppMetrics) Clear() {
    appsMetricsMutex.Lock()
    defer appsMetricsMutex.Unlock()
    (*a) = make(map[AppName]map[Version]map[string]AppStats)
}


func (m AppMetrics) ConvertJsonFriendly() AppMetricsJson {
    appsMetricsMutex.Lock()
    defer appsMetricsMutex.Unlock()
    res := AppMetricsJson{}
    for appName, obj := range m {
        res[appName] = make(map[string]map[string]AppStats)
        for appVersion, appMetrics := range obj {
            res[appName][strconv.Itoa(int(appVersion))] = appMetrics
        }
    }
    return res
}


const (
    CAUTION_INTERVAL = 2
    MINUTES_DELTA = 15 
)

type Needs int

type MemoryNeeds Needs
type CpuNeeds Needs
type NetworkNeeds Needs

type AppNeeds struct {
    MemoryNeeds MemoryNeeds
    CpuNeeds CpuNeeds
    NetworkNeeds NetworkNeeds
}

func (d DeploymentCount) Max(current MaxAble, max MaxAble) MaxAble {
    if max == nil {
        return current
    }
    castMax := max.(DeploymentCount)
    castCurrent := current.(DeploymentCount)
    if castCurrent > castMax {
        castMax = castCurrent
    }
    return castMax
}

func (an AppNeeds) Max(current MaxAble, max MaxAble) MaxAble {
    if max == nil {
        return current
    }
    castMax := max.(AppNeeds)
    castCurrent := current.(AppNeeds)

    if castCurrent.CpuNeeds > castMax.CpuNeeds {
        castMax.CpuNeeds = castCurrent.CpuNeeds
    }
    if castCurrent.MemoryNeeds > castMax.MemoryNeeds {
        castMax.MemoryNeeds = castCurrent.MemoryNeeds
    }
    if castCurrent.NetworkNeeds > castMax.NetworkNeeds {
        castMax.NetworkNeeds = castCurrent.NetworkNeeds
    }

    return castMax
}

type Minutes int

type MaxAble interface {
    Max(MaxAble, MaxAble) MaxAble
}

type TimeBased map[Minutes]MaxAble

type WeekdayBased map[time.Weekday]TimeBased

type WeekdayBasedDeploymentCount struct {
    Based WeekdayBased
}

func (w WeekdayBasedDeploymentCount) Get(day time.Weekday, minutes Minutes) DeploymentCount {
    res := w.Based.get(day, minutes)
    if res == nil {
        return DeploymentCount(0)
    }
    return res.(DeploymentCount)
}

func (w *WeekdayBasedDeploymentCount) Set(day time.Weekday, minutes Minutes, ns DeploymentCount) {
    if len(w.Based) == 0 {
        w.Based = make(map[time.Weekday]TimeBased)
    }
    w.Based.set(day, minutes, ns)
}

func (w *WeekdayBasedDeploymentCount) SetFlat(ns DeploymentCount) {
    if len(w.Based) == 0 {
        w.Based = make(map[time.Weekday]TimeBased)
    }
    w.Based.setFlat(ns)
}

type WeekdayBasedAppNeeds struct {
    Based WeekdayBased
}

func (w WeekdayBasedAppNeeds) Get(day time.Weekday, minutes Minutes) AppNeeds {
    res := w.Based.get(day, minutes)
    if res == nil {
        fmt.Println(w.Based)
        return AppNeeds{}
    }
    return res.(AppNeeds)
}

func (w *WeekdayBasedAppNeeds) Set(day time.Weekday, minutes Minutes, ns AppNeeds) {
    if len(w.Based) == 0 {
        w.Based = make(map[time.Weekday]TimeBased)
    }
    w.Based.set(day, minutes, ns)
}

func (w *WeekdayBasedAppNeeds) SetFlat(ns AppNeeds) {
    if len(w.Based) == 0 {
        w.Based = make(map[time.Weekday]TimeBased)
    }
    w.Based.setFlat(ns)
}

func (t TimeBased) Get(minutes Minutes) MaxAble {
    return t[minutes]
}

func (w WeekdayBased) get(day time.Weekday, minutes Minutes) MaxAble {
    var max MaxAble
    for i := (minutes - CAUTION_INTERVAL * MINUTES_DELTA); i <= minutes + CAUTION_INTERVAL * MINUTES_DELTA; i += MINUTES_DELTA {
        if i >= 0 && i < 1440 {
            current := w[day][Minutes(i)]
            if current != nil {
                max = current.Max(current, max)
            }
        } else if i < 0 {
            var dayBefore time.Weekday
            if day == time.Sunday {
                dayBefore = time.Saturday
            } else {
                dayBefore = day - 1
            }
            currentBefore := w[dayBefore][Minutes(1440 + i)]
            if currentBefore != nil {
                max = currentBefore.Max(currentBefore, max)
            }
        } else if i >= 1440 {
            var dayAfter time.Weekday
            if day == time.Saturday {
                dayAfter = time.Sunday
            } else {
                dayAfter = day + 1
            }
            currentAfter := w[dayAfter][Minutes(i - 1440)]
            if currentAfter != nil {
                max = currentAfter.Max(currentAfter, max)
            }
        }
    }
    return max
}

func (w WeekdayBased) set(day time.Weekday, minutes Minutes, ns MaxAble) {
    if _, exists := w[day]; !exists {
        w[day] = make(map[Minutes]MaxAble)
    }
    w[day][minutes] = ns
}

func (w WeekdayBased) setFlat(ns MaxAble) {
    for i := 0; i < 7; i++ {
        w[time.Weekday(i)] = make(map[Minutes]MaxAble)
        for m := 0; m < (24 * 60); m += MINUTES_DELTA {
            w[time.Weekday(i)][Minutes(m)] = ns
        }
    }
}

//get weekday and minutes in MINUTES_DELTA increments. always rounded down
func timeToWeekdayMinutes(t time.Time) (time.Weekday, Minutes) {
    utcTime := t.UTC()
    w := utcTime.Weekday()
    m := utcTime.Hour() * 60 + utcTime.Minute()
    if m % MINUTES_DELTA != 0 {
        m = int(m / MINUTES_DELTA) * MINUTES_DELTA
    }
    return w, Minutes(m)
}

