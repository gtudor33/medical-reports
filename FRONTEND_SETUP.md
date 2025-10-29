# Medical Reports Frontend Setup

## Overview

A modern, clean React-based UI for the Medical Reports System built with:
- **React 18** with Vite for fast development
- **Tailwind CSS** for beautiful, responsive styling
- **React Router** for navigation
- **Axios** for API communication
- **React Hook Form** for form management

## Features

### Authentication
- Login and registration pages
- JWT token-based authentication
- Protected routes

### Dashboard
- Overview of report statistics
- Recent reports list
- Quick access to create new reports

### Reports Management
- List all reports with filtering (Draft, In Review, Finalized)
- Search by patient name, CNP, or report type
- Create new medical reports
- Edit draft reports
- View detailed report information
- Finalize reports (locks them from editing)
- Delete draft reports

### ICD-10 Integration
- Real-time search for ICD-10 diagnosis codes
- Multi-code selection
- Visual display of selected codes

## Running Locally

### Option 1: With Docker (Recommended)

Build and run everything:
```bash
docker-compose up --build
```

The frontend will be available at: **http://localhost:3000**

### Option 2: Development Mode (Without Docker)

1. Make sure the backend API is running on port 8080

2. Navigate to the frontend directory:
```bash
cd frontend
```

3. Install dependencies:
```bash
npm install
```

4. Start the development server:
```bash
npm run dev
```

The frontend will be available at: **http://localhost:3000**

## Environment Variables

Create a `.env` file in the `frontend` directory:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable components
│   │   ├── Layout.jsx    # Main layout with navigation
│   │   └── ICD10Search.jsx
│   ├── contexts/          # React contexts
│   │   └── AuthContext.jsx
│   ├── pages/             # Page components
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   ├── Dashboard.jsx
│   │   ├── ReportsList.jsx
│   │   ├── ReportForm.jsx
│   │   └── ReportDetail.jsx
│   ├── services/          # API services
│   │   └── api.js
│   ├── App.jsx           # Main app with routing
│   ├── index.css         # Global styles
│   └── main.jsx          # Entry point
├── Dockerfile            # Production build
├── nginx.conf            # Nginx configuration
└── vite.config.js        # Vite configuration
```

## Available Pages

### Public Routes
- `/login` - Login page
- `/register` - Registration page

### Protected Routes (requires authentication)
- `/dashboard` - Dashboard overview
- `/reports` - List of all reports
- `/reports/new` - Create new report
- `/reports/:id` - View report details
- `/reports/:id/edit` - Edit draft report

## Building for Production

```bash
cd frontend
npm run build
```

The production-ready files will be in the `dist` directory.

## Design System

The UI uses a clean, professional design with:
- **Primary Color**: Blue (#0ea5e9)
- **Status Colors**:
  - Draft: Yellow
  - In Review: Orange
  - Finalized: Green
- **Typography**: System fonts for readability
- **Responsive**: Works on mobile, tablet, and desktop

## API Integration

All API calls are centralized in `src/services/api.js`:

- **Auth**: Login, Register
- **Reports**: List, Get, Create, Update, Delete, Finalize
- **Reference Data**: ICD-10 code search

The API client automatically:
- Adds JWT tokens to requests
- Handles 401 errors (redirects to login)
- Provides consistent error handling
