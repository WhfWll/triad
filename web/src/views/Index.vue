<template>
    <div class="dashboard">
        <div class="stats-row">
            <div class="stat-card">
                <div class="stat-icon tasks"><i class="el-icon-s-order"></i></div>
                <div class="stat-body">
                    <div class="stat-top">
                        <span class="stat-label">总任务</span>
                        <span class="stat-trend">+{{latestWeekTaskCount}} 本周</span>
                    </div>
                    <span class="stat-value">{{taskCount}}</span>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon targets"><i class="el-icon-monitor"></i></div>
                <div class="stat-body">
                    <div class="stat-top">
                        <span class="stat-label">总目标</span>
                        <span class="stat-trend warn">高危 {{targetriskstat.highCount}}</span>
                    </div>
                    <span class="stat-value">{{targetriskstat.targetRiskCount}}</span>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon vulns"><i class="el-icon-warning"></i></div>
                <div class="stat-body">
                    <div class="stat-top">
                        <span class="stat-label">总漏洞</span>
                        <span class="stat-trend danger">高危 {{highRiskCount}}</span>
                    </div>
                    <span class="stat-value">{{vulTotal}}</span>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon tools"><i class="el-icon-tools"></i></div>
                <div class="stat-body">
                    <div class="stat-top">
                        <span class="stat-label">工具总数</span>
                        <span class="stat-trend">场景 {{tool.taskSceneCount}}</span>
                    </div>
                    <span class="stat-value">{{tool.vulCount + tool.fingerCount + tool.taskSceneCount}}</span>
                </div>
            </div>
        </div>

        <div class="charts-grid">
            <div class="chart-card">
                <div class="chart-hd">
                    <h3>任务统计</h3>
                    <div class="tab-group">
                        <span :class="{on: taskTab === 1}" @click="switchTaskTab(1)">周</span>
                        <span :class="{on: taskTab === 2}" @click="switchTaskTab(2)">月</span>
                        <span :class="{on: taskTab === 3}" @click="switchTaskTab(3)">年</span>
                    </div>
                </div>
                <div class="chart-bd"><div id="taskTj" style="height:200px"></div></div>
            </div>
            <div class="chart-card">
                <div class="chart-hd"><h3>目标风险</h3></div>
                <div class="chart-bd chart-bd-row">
                    <div class="risk-dots">
                        <div class="rd"><i class="dot h"></i><span>高危</span><em>{{targetriskstat.highCount}}</em></div>
                        <div class="rd"><i class="dot m"></i><span>中危</span><em>{{targetriskstat.mediumCount}}</em></div>
                        <div class="rd"><i class="dot l"></i><span>低危</span><em>{{targetriskstat.lowCount}}</em></div>
                        <div class="rd"><i class="dot s"></i><span>安全</span><em>{{targetriskstat.safeCount}}</em></div>
                    </div>
                    <div id="bar" style="flex:1;height:180px"></div>
                </div>
            </div>
        </div>

        <div class="charts-grid">
            <div class="chart-card chart-wide">
                <div class="chart-hd">
                    <h3>漏洞类型</h3>
                    <div class="tab-group">
                        <span :class="{on: vulnTypeTab===1}" @click="switchVulnTypeTab(1)">周</span>
                        <span :class="{on: vulnTypeTab===2}" @click="switchVulnTypeTab(2)">月</span>
                        <span :class="{on: vulnTypeTab===3}" @click="switchVulnTypeTab(3)">年</span>
                    </div>
                </div>
                <div class="chart-bd"><div id="vuln_type_bar" style="height:200px"></div></div>
            </div>
            <div class="chart-card chart-narrow">
                <div class="chart-hd"><h3>漏洞等级</h3></div>
                <div class="chart-bd">
                    <div id="vuln_pie" style="height:150px"></div>
                    <div class="vbar">
                        <div class="vbar-i" v-for="(it,i) in vulExploitImpact" :key="i">
                            <span class="lbl">{{it.label}}</span>
                            <span class="val" :class="it.color">{{it.value}}</span>
                            <span class="pct">{{it.percentage}}%</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="charts-grid">
            <div class="chart-card chart-wide">
                <div class="chart-hd">
                    <h3>漏洞趋势</h3>
                    <div class="tab-group">
                        <span :class="{on: trendTab===1}" @click="switchTrendTab(1)">周</span>
                        <span :class="{on: trendTab===2}" @click="switchTrendTab(2)">月</span>
                        <span :class="{on: trendTab===3}" @click="switchTrendTab(3)">年</span>
                    </div>
                </div>
                <div class="chart-bd"><div id="vuln_line" style="height:180px"></div></div>
            </div>
            <div class="chart-card chart-narrow">
                <div class="chart-hd"><h3>漏洞取证</h3></div>
                <div class="chart-bd">
                    <div class="ev-grid">
                        <div class="ev-i" v-for="ev in evList" :key="ev.label">
                            <div class="ev-i-icon" :style="{background:ev.bg}"><i :class="ev.icon"></i></div>
                            <div class="ev-i-info"><span class="ev-lbl">{{ev.label}}</span><span class="ev-val">{{ev.count}}</span></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="charts-grid">
            <div class="chart-card chart-full">
                <div class="chart-hd"><h3>最新消息</h3></div>
                <div class="chart-bd">
                    <div class="msgs">
                        <div class="msg-i" v-for="(it,i) in msgstat" :key="i">
                            <span class="msg-t">{{it.createTime}}</span>
                            <span class="msg-c">{{it.content}}</span>
                        </div>
                        <div v-if="msgstat.length===0" class="msg-empty">暂无消息</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.dashboard {
  height: 100%; overflow-y: auto; padding: 4px;

  .stats-row {
    display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 16px;
  }

  .stat-card {
    background: #1a1a2e; border: 1px solid rgba(0,212,170,.08); border-radius: 12px; padding: 20px;
    display: flex; align-items: flex-start; gap: 16px; transition: all .25s;
    &:hover { border-color: rgba(0,212,170,.25); transform: translateY(-2px); box-shadow: 0 8px 30px rgba(0,0,0,.3); }
  }

  .stat-icon {
    width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; color: #fff; flex-shrink: 0;
    &.tasks { background: linear-gradient(135deg,#00d4aa,#00b894); }
    &.targets { background: linear-gradient(135deg,#7c3aed,#6d28d9); }
    &.vulns { background: linear-gradient(135deg,#ef4444,#dc2626); }
    &.tools { background: linear-gradient(135deg,#f59e0b,#d97706); }
  }

  .stat-body { flex: 1; }

  .stat-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }

  .stat-label { font-size: 13px; color: rgba(148,163,184,.6); }
  .stat-trend { font-size: 12px; color: rgba(148,163,184,.5); &.warn { color: #f59e0b; } &.danger { color: #ef4444; } }

  .stat-value { display: block; font-size: 28px; font-weight: 700; color: rgba(226,232,240,.9); }

  .charts-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px; }
  .chart-wide { grid-column: 1 / 2; }
  .chart-narrow { grid-column: 2 / 3; }
  .chart-full { grid-column: 1 / -1; }

  .chart-card {
    background: #1a1a2e; border: 1px solid rgba(0,212,170,.08); border-radius: 12px; overflow: hidden;
    &:hover { border-color: rgba(0,212,170,.15); }
  }

  .chart-hd {
    display: flex; justify-content: space-between; align-items: center; padding: 16px 20px;
    border-bottom: 1px solid rgba(0,212,170,.06);
    h3 { margin: 0; font-size: 14px; font-weight: 600; color: rgba(226,232,240,.85); }
  }

  .tab-group {
    display: flex; background: rgba(15,15,26,.6); border-radius: 6px; padding: 2px;
    span {
      display: inline-block; padding: 4px 12px; font-size: 12px; border-radius: 4px;
      color: rgba(148,163,184,.5); cursor: pointer; transition: all .2s;
      &:hover { color: rgba(226,232,240,.7); }
      &.on { background: #00d4aa; color: #0f0f1a; font-weight: 600; }
    }
  }

  .chart-bd { padding: 16px 20px; }
  .chart-bd-row { display: flex; gap: 16px; }

  .risk-dots {
    min-width: 100px;
    .rd {
      display: flex; align-items: center; gap: 6px; margin-bottom: 8px;
      .dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
        &.h { background: #ef4444; } &.m { background: #f59e0b; } &.l { background: #eab308; } &.s { background: #00d4aa; }
      }
      span { font-size: 13px; color: rgba(148,163,184,.6); flex:1; }
      em { font-size: 14px; font-weight: 600; color: rgba(226,232,240,.8); font-style: normal; }
    }
  }

  .vbar { display: flex; gap: 8px; margin-top: 8px; }
  .vbar-i { flex: 1; text-align: center;
    .lbl { display: block; font-size: 12px; color: rgba(148,163,184,.5); margin-bottom: 4px; }
    .val { display: block; font-size: 18px; font-weight: 700; color: rgba(226,232,240,.8); }
    .pct { display: block; font-size: 11px; color: rgba(148,163,184,.4); margin-top: 2px; }
  }

  .ev-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
  .ev-i { display: flex; align-items: center; gap: 10px; padding: 10px; border-radius: 8px; background: rgba(15,15,26,.4); }
  .ev-i-icon { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 16px; color: #fff; flex-shrink: 0; }
  .ev-i-info { flex:1; .ev-lbl { display: block; font-size: 12px; color: rgba(148,163,184,.5); } .ev-val { font-size: 18px; font-weight: 700; color: rgba(226,232,240,.85); } }

  .msgs { max-height: 200px; overflow-y: auto; }
  .msg-i { display: flex; align-items: center; gap: 16px; padding: 8px 0; border-bottom: 1px solid rgba(0,212,170,.04);
    &:last-child { border-bottom: none; }
  }
  .msg-t { font-size: 12px; color: rgba(148,163,184,.4); white-space: nowrap; min-width: 140px; }
  .msg-c { font-size: 13px; color: rgba(148,163,184,.65); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .msg-empty { text-align: center; padding: 20px; color: rgba(148,163,184,.3); }
}
</style>

<script>
import { system } from '@/api/system.js'

export default {
  name: 'Index',
  data() {
    return {
      taskCount: '--', latestWeekTaskCount: '--',
      targetriskstat: { targetRiskCount: '--', highCount: 0, mediumCount: 0, lowCount: 0, safeCount: 0 },
      vulTotal: 0, highRiskCount: 0,
      vulExploitImpact: [],
      vulevidencestat: { fileLeakCount: 0, dbCount: 0, infoLeakCount: 0, loginCredentialsCount: 0, remoteControlCount: 0 },
      tool: { vulCount: 0, fingerCount: 0, taskSceneCount: 0 },
      msgstat: [],
      taskTab: 1, vulnTypeTab: 1, trendTab: 1,
    }
  },
  computed: {
    evList() {
      const s = this.vulevidencestat;
      return [
        { icon: 'el-icon-folder-opened', label: '文件泄露', count: s.fileLeakCount, bg: 'linear-gradient(135deg,#3b82f6,#2563eb)' },
        { icon: 'el-icon-s-data', label: '数据库', count: s.dbCount, bg: 'linear-gradient(135deg,#00d4aa,#00b894)' },
        { icon: 'el-icon-document-copy', label: '信息泄露', count: s.infoLeakCount, bg: 'linear-gradient(135deg,#7c3aed,#6d28d9)' },
        { icon: 'el-icon-user', label: '登录凭证', count: s.loginCredentialsCount, bg: 'linear-gradient(135deg,#f59e0b,#d97706)' },
        { icon: 'el-icon-monitor', label: '远程控制', count: s.remoteControlCount, bg: 'linear-gradient(135deg,#ef4444,#dc2626)' },
      ];
    },
  },
  created() {
    this.getTaskStat();
    this.getTargetRiskStat();
    this.getVulStat();
    this.getToolStat();
    this.getMsgStat();
  },
  methods: {
    switchTaskTab(t) { this.taskTab = t; this.getTaskStat(); },
    switchVulnTypeTab(t) { this.vulnTypeTab = t; this.getVulStat(); },
    switchTrendTab(t) { this.trendTab = t; this.getVulStat(); },
    getTaskStat() {
      system.getTaskInfoStat().then(r => {
        if (r.data.code === 200) { this.taskCount = r.data.data.totalTaskCount; this.latestWeekTaskCount = r.data.data.latestWeekTaskCount; }
      });
    },
    getTargetRiskStat() {
      system.getTargetRiskStat().then(r => {
        if (r.data.code === 200 && r.data.data) {
          const d = r.data.data.riskStatistics;
          this.targetriskstat = { targetRiskCount: r.data.data.targetRiskCount, highCount: d?.highCount || 0, mediumCount: d?.mediumCount || 0, lowCount: d?.lowCount || 0, safeCount: d?.safeCount || 0 };
        }
      });
    },
    getVulStat() {
      system.getTaskVulRiskStat().then(r => {
        if (r.data.code === 200) {
          const d = r.data.data;
          this.vulTotal = d?.totalVulnerabilities || 0;
          this.highRiskCount = d?.highRiskCount || 0;
          this.vulExploitImpact = (d?.riskLevelStatistics || []).map(it => ({
            label: it.riskName, value: it.count, color: this.getRiskColor(it.riskName),
            percentage: this.vulTotal > 0 ? ((it.count / this.vulTotal) * 100).toFixed(1) : 0,
          }));
        }
      });
      system.getVulEvidenceStat().then(r => {
        if (r.data.code === 200) { this.vulevidencestat = r.data.data || this.vulevidencestat; }
      });
    },
    getToolStat() {
      system.getToolInfoStat().then(r => {
        if (r.data.code === 200) { this.tool = r.data.data || this.tool; }
      });
    },
    getMsgStat() {
      system.getMessageStat().then(r => {
        if (r.data.code === 200) { this.msgstat = r.data.data || []; }
      });
    },
    getRiskColor(name) {
      const m = { '致命': 'danger', '高危': 'danger', '中危': 'warning', '低危': 'info', '信息': 'info' };
      return m[name] || 'info';
    },
  },
};
</script>
