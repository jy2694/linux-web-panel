# Linux Web Panel

리눅스 서버 관리를 위한 웹 기반 관리 패널입니다.

## 기술 스택

### Frontend
- SvelteKit
- TypeScript
- Vite

### Backend
- Go
- Gin Framework

## 프로젝트 구조

```
.
├── frontend/          # SvelteKit 프론트엔드
│   ├── src/          # 소스 코드
│   ├── static/       # 정적 파일
│   └── ...
└── backend/          # Go 백엔드
    ├── cmd/          # 메인 애플리케이션
    ├── internal/     # 내부 패키지
    └── pkg/          # 공개 패키지
```

## 시작하기

### 필수 요구사항
- Node.js 18.x 이상
- Go 1.21 이상
- Git

### 프론트엔드 설정
```bash
cd frontend
npm install
npm run dev
```

### 백엔드 설정
```bash
cd backend
go mod download
go run cmd/main.go
```

## 개발 가이드

### 프론트엔드 개발
- `npm run dev`: 개발 서버 실행
- `npm run build`: 프로덕션 빌드
- `npm run preview`: 프로덕션 빌드 미리보기

### 백엔드 개발
- `go run cmd/main.go`: 개발 서버 실행
- `go test ./...`: 테스트 실행
