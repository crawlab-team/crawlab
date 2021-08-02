<template>
  <svg
      class="img-empty"
      version="1.1"
      viewBox="0 0 79 86"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
  >
    <title>Group 2</title>
    <defs>
      <linearGradient
          id="linearGradient-1"
          x1="38.8503086%"
          x2="61.1496914%"
          y1="0%"
          y2="100%"
      >
        <stop offset="0%" stop-color="#FCFCFD"/>
        <stop offset="100%" stop-color="#EEEFF3"/>
      </linearGradient>
      <linearGradient
          id="linearGradient-2"
          x1="0%"
          x2="100%"
          y1="9.5%"
          y2="90.5%"
      >
        <stop offset="0%" stop-color="#FCFCFD"/>
        <stop offset="100%" stop-color="#E9EBEF"/>
      </linearGradient>
      <rect
          id="path-3"
          height="36"
          width="17"
          x="0"
          y="0"
      />
    </defs>
    <g
        id="Illustrations"
        fill="none"
        fill-rule="evenodd"
        stroke="none"
        stroke-width="1"
    >
      <g id="B-type" transform="translate(-1268.000000, -535.000000)">
        <g id="Group-2" transform="translate(1268.000000, 535.000000)">
          <path
              id="Oval-Copy-2"
              d="M39.5,86 C61.3152476,86 79,83.9106622 79,81.3333333 C79,78.7560045 57.3152476,78 35.5,78 C13.6847524,78 0,78.7560045 0,81.3333333 C0,83.9106622 17.6847524,86 39.5,86 Z"
              fill="#F7F8FC"
          />
          <polygon
              id="Rectangle-Copy-14"
              fill="#E5E7E9"
              points="13 58 53 58 42 45 2 45"
              transform="translate(27.500000, 51.500000) scale(1, -1) translate(-27.500000, -51.500000) "
          />
          <g id="Group-Copy"
             transform="translate(34.500000, 31.500000) scale(-1, 1) rotate(-25.000000) translate(-34.500000, -31.500000) translate(7.000000, 10.000000)">
            <polygon
                id="Rectangle-Copy-10"
                fill="#E5E7E9"
                points="2.84078316e-14 3 18 3 23 7 5 7"
                transform="translate(11.500000, 5.000000) scale(1, -1) translate(-11.500000, -5.000000) "
            />
            <polygon id="Rectangle-Copy-11" fill="#EDEEF2" points="-3.69149156e-15 7 38 7 38 43 -3.69149156e-15 43"/>
            <rect
                id="Rectangle-Copy-12"
                fill="url(#linearGradient-1)"
                height="36"
                transform="translate(46.500000, 25.000000) scale(-1, 1) translate(-46.500000, -25.000000) "
                width="17"
                x="38"
                y="7"
            />
            <polygon
                id="Rectangle-Copy-13"
                fill="#F8F9FB"
                points="24 7 41 7 55 -3.63806207e-12 38 -3.63806207e-12"
                transform="translate(39.500000, 3.500000) scale(-1, 1) translate(-39.500000, -3.500000) "
            />
          </g>
          <rect
              id="Rectangle-Copy-15"
              fill="url(#linearGradient-2)"
              height="36"
              width="40"
              x="13"
              y="45"
          />
          <g id="Rectangle-Copy-17" transform="translate(53.000000, 45.000000)">
            <mask id="mask-4" fill="white">
              <use xlink:href="#path-3"/>
            </mask>
            <use
                id="Mask"
                fill="#E0E3E9"
                transform="translate(8.500000, 18.000000) scale(-1, 1) translate(-8.500000, -18.000000) "
                xlink:href="#path-3"
            />
            <polygon
                id="Rectangle-Copy"
                fill="#D5D7DE"
                mask="url(#mask-4)"
                points="7 0 24 0 20 18 -1.70530257e-13 16"
                transform="translate(12.000000, 9.000000) scale(-1, 1) translate(-12.000000, -9.000000) "
            />
          </g>
          <polygon
              id="Rectangle-Copy-18"
              fill="#F8F9FB"
              points="62 45 79 45 70 58 53 58"
              transform="translate(66.000000, 51.500000) scale(-1, 1) translate(-66.000000, -51.500000) "
          />
        </g>
      </g>
    </g>
  </svg>
</template>

<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'ImgEmpty',
  setup() {
    return {};
  },
});
</script>

<style lang="scss" scoped>
.img-empty {
  max-width: 100%;
  max-height: 100%;
}
</style>
