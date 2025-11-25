FROM mcr.microsoft.com/dotnet/runtime:8.0 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY MobilityCore/MobilityCore.sln ./
COPY ["MobilityCore/MobilityCore.Application/MobilityCore.Application.csproj", "MobilityCore/MobilityCore.Application/"]
COPY ["MobilityCore/MobilityCore.Shared/MobilityCore.Shared.csproj", "MobilityCore/MobilityCore.Shared/"]
RUN dotnet restore "MobilityCore/MobilityCore.Application/MobilityCore.Application.csproj"
COPY . .
WORKDIR "/src/MobilityCore/MobilityCore.Application"
RUN dotnet build "./MobilityCore.Application.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "./MobilityCore.Application.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "MobilityCore.Application.dll"]
