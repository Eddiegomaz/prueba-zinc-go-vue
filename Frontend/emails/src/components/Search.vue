<!-- SearchBar.vue -->
<template>
  <div class="flex items-center space-x-2 my-1">
    <input
      v-model="searchQuery"
      type="text"
      placeholder="Buscar..."
      class="border p-2 w-full rounded-lg"
    />
    <button @click="search" class="px-4 py-2 bg-blue-500 text-white rounded">
      Buscar
    </button>
  </div>
</template>

<!-- <script>
export default {
  data() {
    return {
      searchQuery: "",
    };
  },
  methods: {
    search() {
      // Lógica de búsqueda (puedes implementar tu propia lógica aquí)
      console.log("Buscando:", this.searchQuery);
      // Puedes emitir un evento para que el componente padre maneje la búsqueda
      this.$emit("search", this.searchQuery);
    },
  },
};
</script> -->

<script>
export default {
  // Puedes agregar lógica adicional si es necesario
};
</script>
