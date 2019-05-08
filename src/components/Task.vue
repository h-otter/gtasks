<template>
  <!-- <circle r=15 cx=20 cy=20 v-on:click="" v-on:mouseover=""></circle> -->
  <circle 
    :r="radius"
    :cx="(1 + 2 * vxy[0]) * (margin + radius)"
    :cy="(1 + 2 * vxy[1]) * (margin + radius)"
    v-on:mouseover="$emit('mouseOver', key)"
  ></circle>
    <!-- v-on:click="" -->
</template>

<script>
export default {
  name: "task",
  data() {
    return {
      margin: 5,
      radius: 15,
    };
  },
  props: {
    vxy: Array,
  },
};
</script>

<style>
circle{
  fill: #333;
  transition: all 0.4s cubic-bezier(.96,.04,.04,.96);
  stroke: rgb(119, 0, 255);
  stroke-width: 1px;
}
</style>
