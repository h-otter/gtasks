<template>
  <!-- <circle r=15 cx=20 cy=20 v-on:click="" v-on:mouseover=""></circle> -->
  <circle 
    :r="radius"
    :cx="x"
    :cy="y"
    :class="{'selected': isSelected}"
    v-on:mouseover="$emit('mouseOver', key)"
    v-on:click="$emit('click', key)"
    v-bind:style="{'fill': color}"
  ></circle>
</template>

<script>
export default {
  name: "task-circle",
  data() {
    return {};
  },
  props: {
    radius: Number,
    x: Number,
    y: Number,

    isSelected: Boolean,
    color: String,
  },
};
</script>

<style>
circle{
  /* fill: #ff0000; */
  /* fill: #ffff00; */
  fill: #00ffff;
  transition: all 0.4s cubic-bezier(.96,.04,.04,.96);
  stroke: #333;
  stroke-width: 3px;
}
.selected{
  stroke-width: 10px;
}
</style>
