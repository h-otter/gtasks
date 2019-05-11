<template>
<div id="dashboard">
  <v-app id="inspire">
    <v-navigation-drawer
      fixed
      v-model="drawerRight"
      v-if="tasks[viewingTaskID]"
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
        <v-layout justify-left align-top v-on:click.self="selectTask(null)">
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
              v-on:click="selectTask(id)"
              v-bind:isSelected="isSelected && viewingTaskID == id"
              v-bind:color="stateToColor(tasks[id].state)"
            ></task-circle>

            <!-- month line -->
          </svg>
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

    drawerLeft: false,
    drawerRight: true,

    // right drawer
    menu: false,
    viewingTaskID: 1,
    isSelected: false,
  }),
  props: {
    source: String
  },
  methods: {
    selectTask: function (id) {
      if (this.viewingTaskID != id && this.viewingTaskID) {
        var axios = require("axios")

        // PUT: 失敗したら中断
        axios
          .put("http://127.0.0.1:12345/api/task/"+this.viewingTaskID, this.tasks[this.viewingTaskID])
          .then(response => {})

        this.setGraphCoordinate()
      }

      this.isSelected = !!id;
      this.viewingTaskID = this.isSelected ? id : null;
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

      var l = 2 * this.circleRadius + 10
      for (var n in this.tasks) {
        g.setNode(n, { label: this.tasks[n].title, width: l, height: l });

        if (this.tasks[n].dueDate) {
          var d = new Date(this.tasks[n].dueDate)
          g.setParent(n, (d.getMonth() + 1)+"-"+d.getFullYear());
        }

        if (this.tasks[n].dependsOn) {
          for (var e of this.tasks[n].dependsOn) {
            g.setEdge(e, n)
          }
        }
      }

      dagre.layout(g);

      for (var n in this.tasks) {
        this.tasks[n].x = g.node(n).x
        this.tasks[n].y = g.node(n).y
        // console.log("Node " + n + ": " + JSON.stringify(this.tasks[n]));
      }

      this.paths = []
      for (var e of g.edges()) {
        this.paths.push(g.edge(e).points)
        // console.log("Edge " + e.v + " -> " + e.w + ": " + JSON.stringify(g.edge(e)));
      }

      this.graphWidth = g.graph().width
      this.graphHeight = g.graph().height
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

        if (!this.tasks[this.viewingTaskID]) {
          this.viewingTaskID = 0;
        }

        this.setGraphCoordinate()
      })
  },
}
</script>

<style>

</style>
