<template>
  <div>
    <div class="grid grid-cols-10 h-screen gap-2">
      <!-- Primer contenedor -->
      <div class="col-span-6 backgroung-contain overflow-y-auto">
        <!-- BUSCADOR -->
        <div class="flex items-center space-x-2 my-3">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Buscar..."
            class="border p-2 w-full rounded-lg"
            @keydown.enter.prevent="searchResponse"
          />
          <button
            @click="searchResponse"
            class="px-4 py-2 bg-blue-500 text-white rounded"
          >
            Buscar
          </button>
        </div>
        <!-- VISUAL DE TODOS CORREOS -->
        <ul class="flex flex-col gap-1">
          <li
            v-for="(item, index) in items"
            :key="index"
            class="hover:bg-gray-400 hover:text-white"
            @click="() => showEmail(item._source)"
          >
            <div
              class="grid grid-cols-3 gap-1 rounded border border-gray-400 text-xs"
            >
              <div class="text-sm truncate">
                <p><strong> From </strong>:</p>
                <span> {{ item._source.From }}</span>
              </div>
              <div class="text-sm truncate">
                <p><strong> To </strong>:</p>
                <span>{{ item._source.To }}</span>
              </div>
              <div class="text-sm truncate">
                <p><strong> Subject: </strong>:</p>
                <span class="pr-1"> {{ item._source.Subject }}</span>
              </div>
            </div>
          </li>
        </ul>
        <!-- BOTONES DE PAGINACIÓN -->
        <div class="flex flex-row justify-between">
          <button
            @click="previousPage"
            :disabled="numberPage === 1"
            class="px-4 py-2 bg-blue-500 text-white rounded transition duration-300 ease-in-out transform hover:bg-blue-700 active:bg-blue-800 cursor-pointer"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              class="h-6 w-6"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>
          <p>{{ numberPage }} - {{ Math.floor(totalEmails / sizePage) }}</p>
          <button
            @click="nextPage"
            :disabled="numberPage > totalEmails"
            class="px-4 py-2 bg-blue-500 text-white rounded transition duration-300 ease-in-out transform hover:bg-blue-700 active:bg-blue-800 cursor-pointer"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              class="h-6 w-6"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>
        </div>
      </div>
      <!-- Segundo contenedor -->
      <div
        class="col-span-4 flex justify-left items-top border border-gray-200 shadow-2xl overflow-y-auto h-max-96"
        style="max-width: 100%; overflow-x: hidden"
      >
        <div v-if="selectedEmail" class="m-8">
          <p>
            <strong>From:</strong><br />
            {{ selectedEmail.From }}
          </p>
          <p><strong>To:</strong><br />{{ selectedEmail.To }}</p>
          <p><strong>Subject:</strong> <br />{{ selectedEmail.Subject }}</p>
          <p><strong>Content:</strong> <br />{{ selectedEmail.Content }}</p>
        </div>

        <img
          v-else
          src="/welcome.png"
          alt="image thing"
          class="object-contain max-w-full h-auto"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";

let filter = false;
let numberPage = 1;
let sizePage = 15;
let items = ref([]);
let totalEmails = ref([]);
let selectedEmail = ref(null);
let searchQuery = ref("");

// Método para cargar la información

onMounted(() => {
  loadData("", numberPage);
});

async function loadData(data, startQuery) {
  try {
    const response = await axios.get("http://localhost:3000/query", {
      params: {
        query: data,
        from: startQuery,
        size: sizePage,
      },
    });

    // Almacena los resultados en itemsArray
    items.value = response.data.Hits.Hits;
    totalEmails.value = response.data.Hits.Total.Value;
  } catch (error) {
    console.error("Error al cargar la información:", error.message);
  }
}

function showEmail(email) {
  selectedEmail.value = email;
}
function searchResponse() {
  filter = true;
  numberPage = 1;
  loadData(searchQuery.value, numberPage);
}

function previousPage() {
  if (numberPage >= 1) {
    numberPage -= 1;

    loadData(searchQuery.value, numberPage);
  }
}

function nextPage() {
  console.log(numberPage);
  console.log(totalEmails.value);
  if (numberPage < Math.floor(totalEmails.value / sizePage)) {
    numberPage += 1;

    loadData(searchQuery.value, numberPage);
  }
}
</script>

<style scoped>
.backgroung-contain {
  max-width: 100%;
  overflow-x: hidden;
  background-color: #ffffff;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='20' height='20' viewBox='0 0 100 100'%3E%3Cg stroke='%23CCC' stroke-width='0' %3E%3Crect fill='%23F5F5F5' x='-60' y='-60' width='65' height='240'/%3E%3C/g%3E%3C/svg%3E");
}
</style>
