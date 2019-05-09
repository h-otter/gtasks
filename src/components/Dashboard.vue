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
          ></v-text-field>

          <!-- state TOOD: default valueがない -->
          <v-select
            v-model="tasks[viewingTaskID].state"
            :items="status"
            label="State"
          ></v-select>

          <v-textarea
            v-model="tasks[viewingTaskID].description"
            label="Description"
            auto-grow
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
          >
            <template v-slot:activator="{ on }">
              <v-text-field
                v-model="tasks[viewingTaskID].dueDate"
                label="Due date"
                clearable
                readonly
                v-on="on"
              ></v-text-field>
            </template>
            <v-date-picker v-model="tasks[viewingTaskID].dueDate" no-title scrollable>
              <v-spacer></v-spacer>
              <v-btn flat color="primary" @click="menu = false">Cancel</v-btn>
              <v-btn flat color="primary" @click="$refs.menu.save(tasks[viewingTaskID].dueDate)">OK</v-btn>
            </v-date-picker>
          </v-menu>

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
      <!-- <v-toolbar-side-icon @click.stop="drawerLeft = !drawerLeft"></v-toolbar-side-icon> -->
      <!-- <v-toolbar-title>Toolbar</v-toolbar-title> -->
      <v-spacer></v-spacer>
      <!-- plus: add task -->
      <v-toolbar-side-icon @click.stop="drawerRight = !drawerRight"></v-toolbar-side-icon>
    </v-toolbar>

    <!-- <v-navigation-drawer
      fixed
      v-model="drawerLeft"
      app
    >
      <v-list dense>
        <v-list-tile @click.stop="left = !left">
          <v-list-tile-action>
            <v-icon>exit_to_app</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Open Temporary Drawer</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer> -->

    <v-content>
      <v-container fluid fill-height>
        <v-layout justify-left align-top>
          <svg :width="this.graphWidth" :height="this.graphHeight">
            <!-- tasks -->
            <task-circle v-for="(value, id) in tasks"
              :key="id"
              :radius="circleRadius"
              v-bind:x="tasks[id].x"
              v-bind:y="tasks[id].y"
              v-on:mouseOver="updateViewingTask(id)"
              v-on:click="toggleIsSelected(id)"
              v-bind:isSelected="isSelected && viewingTaskID == id"
              v-bind:color="stateToColor(tasks[id].state)"
            ></task-circle>
            <!-- dependency line -->
            <task-path v-for="(points, index) in paths"
              :key="index"
              v-bind:points="points"
            ></task-path>

            <!-- month line -->
          </svg>

<!-- 
          <v-flex shrink>
            <v-tooltip right>
              <v-btn
                icon
                large
                :href="source"
                target="_blank"
                slot="activator"
              >
                <v-icon large>code</v-icon>
              </v-btn>
              <span>Source</span>
            </v-tooltip>
          </v-flex> -->
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
    tasks: {
      1: {
        title: "testing",
        state: "Doing",
        labels: {
          "test-label": "hoge",
        },
        description: "foobar",
        dueDate: new Date().toISOString().substr(0, 10),
        dependsOn: [
          2,
        ],
      },
      3: {
        title: "ToDo",
        state: "ToDo",
        labels: {
          "test-label": "hoge",
        },
        description: "foobar",
        dueDate: new Date().toISOString().substr(0, 10),
        dependsOn: [],
      },
      2: {
        title: "changed",
        state: "Done",
        labels: {
          "test-label": "hoge",
        },
        description: "foobar",
        dueDate: new Date().toISOString().substr(0, 10),
        dependsOn: [],
      },
    },
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
    updateViewingTask: function (id) {
      if (!this.isSelected) {
        this.viewingTaskID = id
      }
    },
    toggleIsSelected: function (id) {
      this.isSelected = !this.isSelected
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
    setGraphCoordinate: function() {
      var dagre = require("dagre");
      var g = new dagre.graphlib.Graph();
      g.setGraph({});
      g.setDefaultEdgeLabel(function() { return {}; });

      // g.setParent(month+"-"+category, month);

      var l = 2 * this.circleRadius + 10
      for (var n in this.tasks) {
        g.setNode(n, { label: this.tasks[n].title, width: l, height: l });

        // g.setParent(n, month+"-"+category);

        for (var e of this.tasks[n].dependsOn) {
          g.setEdge(e, n)
        }
      }

      dagre.layout(g);

      for (var n of g.nodes()) {
        this.tasks[n].x = g.node(n).x
        this.tasks[n].y = g.node(n).y
        console.log("Node " + n + ": " + JSON.stringify(this.tasks[n]));
      }

      this.paths = []
      for (var e of g.edges()) {
        this.paths.push(g.edge(e).points)
        console.log("Edge " + e.v + " -> " + e.w + ": " + JSON.stringify(g.edge(e)));
      }

      this.graphWidth = g.graph().width
      this.graphHeight = g.graph().height
    }
  },
  mounted: function() {
    this.setGraphCoordinate()
  },
}
</script>

<style>

</style>
