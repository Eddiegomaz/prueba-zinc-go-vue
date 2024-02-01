<template>
  <div>
    <ul class="flex flex-col gap-1">
      <li
        v-for="(item, index) in items"
        :key="index"
        class="group hover:bg-white hover:shadow-lg hover:border-transparent transform motion-safe:hover:scale-105"
      >
        <div
          class="grid grid-cols-3 gap-1 rounded border border-gray-400 text-xs"
        >
          <div class="text-sm">
            <p><strong> From </strong>:</p>
            <span> {{ item._source.From }}</span>
          </div>
          <div class="text-sm">
            <p><strong> To </strong>:</p>
            <span>{{ item._source.To }}</span>
          </div>
          <div class="pr-1">
            <span> {{ truncate(item._source.Subject) }}</span>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "ListGroup",
  props: {
    items: {
      type: Object,
      required: true,
    },
  },

  methods: {
    truncate(subject) {
      return subject.length > 30 ? subject.substring(0, 50) + "..." : subject;
    },
  },
};

// console.log(items);
</script>
