<template>
  <!-- <circle r=15 cx=20 cy=20 v-on:click="" v-on:mouseover=""></circle> -->
  <circle 
    :r="radius"
    :cx="(1 + 2 * vxy[0]) * (margin + radius)"
    :cy="(1 + 2 * vxy[1]) * (margin + radius)"
    :class="{'selected': isSelected}"
    v-on:mouseover="$emit('mouseOver', key)"
    v-on:click="$emit('click', key)"
  ></circle>
</template>

<script>
export default {
  name: "task",
  data() {
    return {
      margin: 5,
      radius: 15,
    // state: "doing",

    };
  },
  props: {
    // key: Number,
    // やはりvxyをうける必要はない、このコンポーネントからロジックは切り離す
    vxy: Array,

    isSelected: Boolean,
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
