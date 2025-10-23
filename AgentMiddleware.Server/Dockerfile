FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Copy solution and project files
COPY ["AgentMiddleware.Server/AgentMiddleware.Server.csproj", "AgentMiddleware.Server/"]
RUN dotnet restore "AgentMiddleware.Server/AgentMiddleware.Server.csproj"

# Copy source code
COPY . .
WORKDIR "/src/AgentMiddleware.Server"
RUN dotnet build "AgentMiddleware.Server.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "AgentMiddleware.Server.csproj" -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app
COPY --from=publish /app/publish .

EXPOSE 8080 8081
ENTRYPOINT ["dotnet", "AgentMiddleware.Server.dll"]
