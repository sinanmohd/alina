<template>
  <main class="flex flex-col min-h-screen">
    <header>
      <nav class="max-w-7xl px-3 mx-auto py-3 flex justify-between">
        <section class="flex gap-1 text-3xl">
          <Icon name="my-icon:logo" />
          <p class="font-bold text-rose-700">alina</p>
        </section>

        <section class="flex">
          <Button variant="ghost">
            <a href="https://github.com/sinanmohd/alina" target="_blank">
              <Icon name="uil:github" class="text-2xl"/>
            </a>
          </Button>
        </section>
      </nav>
    </header>

    <Separator />

    <div class="flex-grow">
      <slot />
    </div>

    <Separator />

    <footer>
      <div class="max-w-7xl px-3 mx-auto py-2 flex justify-between">
        <p class="text-sm text-(--ui-text-muted) my-auto">
          Copyright © {{ new Date().getFullYear() }} Sinan -
          <a href="https://www.gnu.org/licenses/agpl-3.0.en.html" target="_blank" class="hover:underline"> AGPL3 License </a>
        </p>

        <section class="flex">
          <Button variant="ghost">
            <a href="https://www.sinanmohd.com/" target="_blank">
              <Icon name="mdi:web" class="text-2xl"/>
            </a>
          </Button>
          <Button variant="ghost">
            <a href="https://www.linkedin.com/in/sinanmohd/" target="_blank">
              <Icon name="uil:linkedin" class="text-2xl"/>
            </a>
          </Button>
          <Button variant="ghost">
            <a href="https://github.com/sinanmohd" target="_blank">
              <Icon name="uil:github" class="text-2xl"/>
            </a>
          </Button>
        </section>
      </div>
    </footer>
  </main>
</template>
