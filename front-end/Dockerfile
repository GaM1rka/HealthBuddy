# Этап сборки Flutter-приложения
FROM ghcr.io/cirruslabs/flutter:stable AS build

WORKDIR /app

# Копируем только необходимые файлы
COPY pubspec.yaml ./
COPY lib/ lib/
COPY web/ web/

# Устанавливаем зависимости и собираем проект
RUN flutter pub get
RUN flutter build web --release

# Этап развертывания
FROM nginx:stable-alpine
COPY --from=build /app/build/web /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]