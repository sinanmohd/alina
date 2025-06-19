<div align = center>

# alina

[![Badge Matrix]](https://matrix.to/#/#chat:sinanmohd.com)
![Badge Pull Requests]
![Badge Issues]

Your awesome frenly neighbourhood file sharing website. 

</div>

# features

- Pause, resume, and retry uploads.
- Prometheus/Grafana metrics.
- Upload files directly from the terminal.
- Parallel chunked uploads.
- Secure chunked uploads without requiring a login.
- Auto Merge duplicate files to save space.
- Markdown rendering for text notes.
- IP based rate limiting.
- Sleek, design is very humane :P.
- and much more...

# gallery

![Preview A]

# build
```
nix develop

cd ./frontend
pnpm i
pnpm exec nuxt generate
cd -

rm -rf ./backend/internal/server/frontend
cp -rv ./frontend/.output/public ./backend/internal/server/frontend

cd ./backend
go build -o ../alina ./cmd/alina/main.go
cd -
```

# development
```
# frontend
nix develop
cd ./frontend
pnpm install
NUXT_PUBLIC_SERVER_URL=http://localhost:8008 pnpm exec nuxt dev --host

# backend
nix develop
mkdir -p ./backend/internal/server/frontend
touch ./backend/internal/server/frontend/stub
cd ./backend
air
```

# special thanks

**[Tailwind, Nuxt, Shadcn]** - *For helping with frontend*

**[Nix, Go]** - *For being cool*

**[Glass Shelf]** - *For breaking my arm so i have free time to do this*


<!----------------------------------{ Thanks }--------------------------------->

[Tailwind, Nuxt, Shadcn]: https://tailwindcss.com/
[Nix, Go]: https://nixos.org/
[Glass Shelf]: https://www.amazon.com/SAYAYO-Floating-Shelves-Tempered-Bathroom/dp/B0CGXB13CR

<!----------------------------------{ Images }--------------------------------->

[Preview A]: https://static.sinanmohd.com/git/alina.png

<!----------------------------------{ Badges }--------------------------------->

[Badge Matrix]: https://img.shields.io/matrix/chat:sinanmohd.com.svg?label=%23chat%3Asinanmohd.com&logo=matrix&server_fqdn=sinanmohd.com
[Badge Issues]: https://img.shields.io/github/issues/sinanmohd/alina
[Badge Pull Requests]: https://img.shields.io/github/issues-pr/sinanmohd/alina
