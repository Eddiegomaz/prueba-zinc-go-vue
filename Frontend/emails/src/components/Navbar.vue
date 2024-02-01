<!-- Navbar.vue -->
<template>
  <nav class="bg-white border-gray-200 dark:border-gray-600 dark:bg-gray-900">
    <div
      class="flex flex-wrap justify-between items-center mx-auto max-w-screen-xl p-4"
    >
      <img src="/donut-chart.png" class="h-8" alt="Flowbite Logo" />
      <span class="self-center text-2xl font-semibold whitespace-nowrap"
        >IndexEmail</span
      >
    </div>
    <div
      class="mt-1 border-gray-200 shadow-sm bg-gray-50 md:bg-white border-y dark:bg-gray-800 dark:border-gray-600"
    ></div>
  </nav>
</template>

<script>
export default {
  name: "Navbar",
};
</script>

<style scoped>
/* Estilos específicos de la barra de navegación */

ul {
  list-style-type: none;
  margin: 0;
  padding: 0;
  display: flex;
}

li {
  margin-right: 10px;
}

a {
  color: white;
  text-decoration: none;
}
</style>
