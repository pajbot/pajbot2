FROM golang:stretch AS build
RUN apt-get update
RUN apt-get install apt-transport-https -y
RUN wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.asc.gpg
RUN mv microsoft.asc.gpg /etc/apt/trusted.gpg.d/
RUN wget -q https://packages.microsoft.com/config/debian/9/prod.list
RUN mv prod.list /etc/apt/sources.list.d/microsoft-prod.list
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get update
RUN apt-get install dotnet-sdk-2.2=2.2.100-1 nodejs -y
ADD . /src
RUN cd /src && ./utils/install.sh
RUN cd /src/web && npm i && npm run build
RUN cd /src/cmd/bot && go build -v -tags csharp

FROM mcr.microsoft.com/dotnet/core/runtime:2.2-bionic
WORKDIR /app
ENV LIBCOREFOLDER /usr/share/dotnet/shared/Microsoft.NETCore.App/2.2.6
ENV PAJBOT2_WEB_PATH /app/web/
COPY --from=build /src/web/static /app/web/static
COPY --from=build /src/web/views /app/web/views
COPY --from=build /src/cmd/bot/bot /app/
COPY --from=build /src/migrations /app/migrations/
COPY --from=build /src/cmd/bot/*.dll /app/
COPY --from=build /src/cmd/bot/charmap.bin.gz /app/
RUN chmod 777 /app/charmap.bin.gz
CMD ["./bot"]
