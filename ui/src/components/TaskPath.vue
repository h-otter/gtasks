<template>
  <!-- <circle r=15 cx=20 cy=20 v-on:click="" v-on:mouseover=""></circle> -->
  <path
    :d="pointsToD()"
  ></path>
</template>

<script>
export default {
  name: "task-path",
  data() {
    return {};
  },
  props: {
    points: Array,
  },
  methods: {
    pointsToD: function () {
      return "M" + this.points.map(function(e) {
        return e.x + "," + e.y;
      }).join("L")
    },
  }
};
</script>

<style>
path {
  fill: none;
  stroke: #333;
  stroke-width: 2px;
}
</style>
