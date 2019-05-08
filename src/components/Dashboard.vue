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

    <v-navigation-drawer
      temporary
      v-model="left"
      fixed
    ></v-navigation-drawer>

    <v-content>
      <v-container fluid fill-height>
        <v-layout justify-left align-top>
          <svg width="300" height="300">
            <!-- <task></task> -->
            <task v-for="(value, id) in tasks"
              :key="id"
              :vxy="[id, 0]"
              v-on:mouseOver="updateViewingTask(id)"
            ></task>
            <!-- month line -->
            <!-- dependency line -->
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

    <v-navigation-drawer
      right
      temporary
      v-model="right"
      fixed
    ></v-navigation-drawer>
  </v-app>
</div>
</template>

<script>
import Task from "./Task";

export default {
  name: 'Dashboard',
  components: {
    Task,
  },
  data: () => ({
    tasks: {
      1: {
        title: "testing",
        labels: {
          "test-label": "hoge",
        },
        description: "foobar",
        dueDate: new Date().toISOString().substr(0, 10),
        dependsOn: [],
      },
      2: {
        title: "changed",
        labels: {
          "test-label": "hoge",
        },
        description: "foobar",
        dueDate: new Date().toISOString().substr(0, 10),
        dependsOn: [],
      },
    },
    vCoordinate: {},

    drawerLeft: false,
    drawerRight: true,
    right: null,
    left: null,

    // right drawer
    menu: false,
    viewingTaskID: 1,
  }),

  props: {
    source: String
  },
  methods: {
    updateViewingTask: function (id) {
      this.viewingTaskID = id
    },
  }
}
</script>

<style>

</style>
