<template>
<div id="dashboard">
  <v-app id="inspire">
    <v-navigation-drawer
      fixed
      v-model="drawerRight"
      right
      clipped
      app
    >
      <v-form>
        <v-container>
          <v-text-field
            v-model="tasks[viewingTaskID].title"
            label="Title"
            required
            :disabled="!isSelected"
          ></v-text-field>

          <!-- state TOOD: default valueがない -->
          <v-select
            v-model="tasks[viewingTaskID].state"
            :items="status"
            label="State"
            :disabled="!isSelected"
          ></v-select>

          <v-textarea
            v-model="tasks[viewingTaskID].description"
            label="Description"
            auto-grow
            :disabled="!isSelected"
          ></v-textarea>
          <v-menu
            ref="menu"
            v-model="menu"
            :close-on-content-click="false"
            :nudge-right="40"
            :return-value.sync="tasks[viewingTaskID].dueDate"
            lazy
            transition="scale-transition"
            offset-y
            full-width
            min-width="290px"
            :disabled="!isSelected"
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="tasks[viewingTaskID].dueDate"
                label="Due date"
                clearable
                readonly
                :disabled="!isSelected"
                v-on="on"
              ></v-text-field>
            </template>
            <v-date-picker v-model="tasks[viewingTaskID].dueDate" no-title scrollable>
              <v-spacer></v-spacer>
              <v-btn flat color="primary" @click="menu = false">Cancel</v-btn>
              <v-btn flat color="primary" @click="$refs.menu.save(tasks[viewingTaskID].dueDate)">OK</v-btn>
            </v-date-picker>
          </v-menu>

          <!-- category -->
          <!-- labels -->
          <!-- depends on -->
        </v-container>
      </v-form>
    </v-navigation-drawer>

    <v-toolbar
      color="blue-grey"
      dark
      fixed
      app
      clipped-right
    >
      <v-spacer></v-spacer>
      <!-- plus: add task -->
      <v-toolbar-side-icon @click.stop="drawerRight = !drawerRight"></v-toolbar-side-icon>
    </v-toolbar>

    <v-content>
      <v-container fluid fill-height>
        <v-layout class="main-layout" justify-left align-top>
          <div class="main-content" id="background" :style="{ gridTemplateRows: gridTemplateRowsJoined }">
            <template v-for="(value, index) in graphHeightByMonth">
              <div class="head" :class="[ index % 2 == 1 ? 'odd' : 'even']" :style="{ gridRow: (index + 1) + '/ span 1', gridColumn: '1 / span 1'}">{{ value.year }}/{{ value.month }}</div>
              <div class="content" :class="[ index % 2 == 1 ? 'odd' : 'even']" :style="{ gridRow: (index + 1) + '/ span 1', gridColumn: '2 / span 1'}"></div>
            </template>
          </div>
          <div class="main-content" id="graph">
            <svg :width="this.graphWidth" :height="this.graphHeight">
              <!-- dependency line -->
              <task-path v-for="(points, index) in paths"
                :key="'task-path-'+index"
                v-bind:points="points"
              ></task-path>
              <!-- tasks -->
              <task-circle v-for="(value, id) in tasks"
                :key="'task-circle-'+id"
                :radius="circleRadius"
                v-bind:x="tasks[id].x"
                v-bind:y="tasks[id].y"
                v-on:mouseOver="updateViewingTask(id)"
                v-on:click="toggleIsSelected(id)"
                v-bind:isSelected="isSelected && viewingTaskID == id"
                v-bind:color="stateToColor(tasks[id].state)"
              ></task-circle>

              <!-- month line -->
            </svg>
          </div>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</div>
</template>

<script>
import TaskCircle from "./TaskCircle";
import TaskPath from "./TaskPath";

export default {
  name: 'Dashboard',
  components: {
    TaskCircle,
    TaskPath,
  },
  data: () => ({
    status: ["ToDo", "Doing", "Done"],
    tasks: {},
    paths: [],

    circleMargin: 5,
    circleRadius: 15,
    graphWidth: 100,
    graphHeight: 100,
    graphHeightByMonth: [],

    drawerLeft: false,
    drawerRight: true,

    // right drawer
    menu: false,
    viewingTaskID: 1,
    isSelected: false,
  }),
  computed: {
    gridTemplateRowsJoined: function() {
      return this.graphHeightByMonth.map(x => x.height + "px").join(" ");
    }
  },
  props: {
    source: String
  },
  methods: {
    updateViewingTask: function (id) {
      if (!this.isSelected) {
        this.viewingTaskID = id
      }
    },
    toggleIsSelected: function (id) {
      if (this.isSelected) {
        var axios = require("axios")

        // PUT: 失敗したら中断
        axios
          .put("http://127.0.0.1:12345/api/task/"+id, this.tasks[id])
          .then(response => {})

        this.setGraphCoordinate()
      }

      this.isSelected = !this.isSelected

      if (this.isSelected) {
        this.viewingTaskID = id
      }
    },
    stateToColor: function(state) {
      switch (state) {
      case "ToDo":
        return "#ffff00"
      case "Doing":
        return "#ff0000"
      case "Done":
        return "#00ffff"
      }
    },
    // addDependency: function() {
      // check this is DAG
      // try {
      //   console.log(dagre.graphlib.alg.topsort(g))
      // } catch(error) {}
    // },
    setGraphCoordinate: function() {
      var dagre = require("dagre");
      var g = new dagre.graphlib.Graph({ compound: true });
      g.setGraph({});
      g.setDefaultEdgeLabel(function() { return {}; });

      // g.setParent(month+"-"+category, month);

      const groupNameOfTask = (task) => {
        const d = new Date(task.dueDate);
        const name = `${d.getMonth() + 1}-${d.getFullYear()}`;
        g.setNode(name, { label: name, width: 500, height: 100 });

        return name;
      }

      const l = 2 * this.circleRadius + 10
      for (const n in this.tasks) {
        const task = this.tasks[n];
        g.setNode(n, { label: task.title, width: l, height: l });

        if (task.dueDate) {
          g.setParent(n, groupNameOfTask(task));
        }

        if (task.dependsOn) {
          for (var e of this.tasks[n].dependsOn) {
            g.setEdge(e, n)
          }
        }
      }

      dagre.layout(g);

      for (let n in this.tasks) {
        this.tasks[n].x = g.node(n).x
        this.tasks[n].y = g.node(n).y
        // console.log("Node " + n + ": " + JSON.stringify(this.tasks[n]));
      }

      this.paths = []
      for (const e of g.edges()) {
        this.paths.push(g.edge(e).points)
        // console.log("Edge " + e.v + " -> " + e.w + ": " + JSON.stringify(g.edge(e)));
      }

      this.graphWidth = g.graph().width
      this.graphHeight = g.graph().height

      this.updateGraphHeightByMonth();
    },
    updateGraphHeightByMonth: function() {
      let minDate = null;
      let maxDate = null;
      for (const n in this.tasks) {
        const task = this.tasks[n];
        if (!task.dueDate) { continue; }

        const d = new Date(task.dueDate);
        const fd = new Date(d.getFullYear(), d.getMonth(), 1); // first day of month

        if (!minDate) {
          minDate = fd;
          maxDate = fd;
        }
        if (fd < minDate) {
          minDate = fd;
        } else if (maxDate < fd) {
          maxDate = fd;
        }
      }

      if (!minDate) {
        // dueDate is not set in all tasks
        this.graphHeightByMonth = [];
      }

      const getFirstDayOfNextMonth = (d) => {
        if (d.getMonth() == 11) {
          return new Date(d.getFullYear() + 1, 0, 1);
        }
        return new Date(d.getFullYear(), d.getMonth() + 1, 1);
      }

      let d = minDate;
      let graphHeightByMonth = []
      do {
        // TODO: get appropriate height
        console.log({ year: d.getFullYear(), month: d.getMonth() + 1, height: 100})
        graphHeightByMonth.push({ year: d.getFullYear(), month: d.getMonth() + 1, height: 100});
        d = getFirstDayOfNextMonth(d);
      } while (d.getTime() <= maxDate.getTime());
      this.graphHeightByMonth = graphHeightByMonth
    }
  },
  mounted: function() {
    // list all
    var axios = require("axios")
    axios
      .get("http://127.0.0.1:12345/api/task/")
      .then(response => {
        this.tasks = {}
        for (var r of response.data) {
          this.tasks[r.id] = r
        }

        this.setGraphCoordinate()
      })
  },
}
</script>

<style>
.main-layout {
  position: relative;
}

.main-content {
  position: absolute;
  display: grid;
  grid-template-columns: 150px 1fr;
  width: 100%;
}

.main-content#background .odd {
  border-top: 1px solid rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(0, 0, 0, 0.2);
  background: rgba(0, 0, 0, 0.1);
}

.main-content#background .head {
  border-right: 1px dashed rgba(0, 0, 0, 0.2);
}

.main-content#graph {
  grid-template-rows: 1fr;
  grid-template-areas: "unused graph";
}

.main-content#graph > svg {
  grid-area: graph;
}
</style>
