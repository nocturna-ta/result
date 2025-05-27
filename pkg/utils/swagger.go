package utils

const ClaudeDarkTheme = `
        /* Claude Dark Theme for Swagger UI */
        body {
            background-color: #1a1a24;
            color: #e0e0e8;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
        }
        
        .swagger-ui .topbar {
            background-color: #2d2b55;
            border-bottom: 2px solid #6e56cf;
        }
        
        .swagger-ui .topbar .download-url-wrapper .select-label {
            color: #e0e0e8;
        }
        
        .swagger-ui .topbar .download-url-wrapper input[type=text] {
            background-color: #312e5c;
            color: #e0e0e8;
            border: 1px solid #6e56cf;
        }
        
        .swagger-ui .info {
            background-color: #252436;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        
        .swagger-ui .info .title {
            color: #a393f0;
            font-weight: 600;
        }
        
        .swagger-ui .info a {
            color: #a393f0;
        }
        
        .swagger-ui .info a:hover {
            color: #c4b5fa;
        }
        
        .swagger-ui .opblock-tag {
            border-radius: 8px;
            background-color: #252436;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
            color: #e0e0e8;
            border: none;
        }
        
        .swagger-ui .opblock-tag:hover {
            background-color: #312e5c;
            color: #fff;
        }
        
        .swagger-ui .opblock {
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
            border: none;
            margin-bottom: 15px;
            background-color: #252436;
        }
        
        .swagger-ui .opblock .opblock-summary {
            border-radius: 8px 8px 0 0;
        }
        
        .swagger-ui .opblock .opblock-summary-path {
            color: #e0e0e8;
        }
        
        .swagger-ui .opblock .opblock-summary-description {
            color: #a9a9b8;
        }
        
        .swagger-ui .opblock-summary-method {
            border-radius: 4px;
            font-weight: 600;
        }
        
        /* POST method color */
        .swagger-ui .opblock-post {
            background-color: rgba(110, 86, 207, 0.1);
            border-color: #6e56cf;
        }
        
        .swagger-ui .opblock-post .opblock-summary-method {
            background-color: #6e56cf;
        }
        
        /* GET method color */
        .swagger-ui .opblock-get {
            background-color: rgba(97, 175, 254, 0.1);
            border-color: #61affe;
        }
        
        /* PUT method color */
        .swagger-ui .opblock-put {
            background-color: rgba(252, 161, 83, 0.1);
            border-color: #fca153;
        }
        
        /* DELETE method color */
        .swagger-ui .opblock-delete {
            background-color: rgba(249, 62, 62, 0.1);
            border-color: #f93e3e;
        }
        
        .swagger-ui .opblock-body {
            background-color: #1e1d2d;
        }
        
        .swagger-ui .btn {
            border-radius: 6px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
        
        .swagger-ui .btn.execute {
            background-color: #6e56cf;
            color: white;
            border-color: #6e56cf;
        }
        
        .swagger-ui .btn.execute:hover {
            background-color: #8570d8;
        }
        
        .swagger-ui .btn.authorize {
            background-color: #6e56cf;
            color: white;
            border-color: #6e56cf;
        }
        
        .swagger-ui .btn.authorize:hover {
            background-color: #8570d8;
        }
        
        .swagger-ui section.models {
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
            background-color: #252436;
            border: none;
        }
        
        .swagger-ui section.models.is-open h4 {
            border-bottom-color: #3d3a63;
            color: #e0e0e8;
        }
        
        .swagger-ui .model-container {
            background-color: #1e1d2d;
            border-radius: 4px;
        }
        
        .swagger-ui .model-box {
            background-color: #252436;
        }
        
        .swagger-ui .model {
            color: #e0e0e8;
        }
        
        .swagger-ui .model-title {
            color: #a393f0;
        }
        
        .swagger-ui input[type=text], 
        .swagger-ui textarea {
            border-radius: 6px;
            border: 1px solid #3d3a63;
            background-color: #1e1d2d;
            color: #e0e0e8;
        }
        
        .swagger-ui select {
            border-radius: 6px;
            border: 1px solid #3d3a63;
            background-color: #1e1d2d;
            color: #e0e0e8;
        }
        
        .swagger-ui select option {
            background-color: #1e1d2d;
        }
        
        .swagger-ui .parameter__name {
            color: #a393f0;
            font-weight: 600;
        }
        
        .swagger-ui .parameter__type {
            color: #a9a9b8;
        }
        
        .swagger-ui table {
            border-radius: 8px;
            overflow: hidden;
            background-color: #1e1d2d;
        }
        
        .swagger-ui table thead tr th {
            background-color: #252436;
            color: #a393f0;
            border-color: #3d3a63;
        }
        
        .swagger-ui table tbody tr td {
            border-color: #3d3a63;
            color: #e0e0e8;
        }
        
        /* Response section styling */
        .swagger-ui .responses-table .response {
            border-radius: 4px;
        }
        
        .swagger-ui .response-col_status {
            color: #a393f0;
        }
        
        .swagger-ui .prop-type {
            color: #a393f0;
        }
        
        .swagger-ui .prop-format {
            color: #a9a9b8;
        }
        
        /* Toggle button colors */
        .swagger-ui .model-box-control:focus, 
        .swagger-ui .models-control:focus, 
        .swagger-ui .opblock-summary-control:focus {
            outline: none;
        }
        
        /* JSON/Schema styling */
        .swagger-ui .microlight {
            background-color: #1e1d2d;
            color: #e0e0e8;
            border-radius: 4px;
            padding: 8px;
        }
        
        /* Code highlighting */
        .swagger-ui .microlight .pun,
        .swagger-ui .microlight .opn, 
        .swagger-ui .microlight .clo {
            color: #e0e0e8;
        }
        
        .swagger-ui .microlight .str {
            color: #a8ff60;
        }
        
        .swagger-ui .microlight .kwd {
            color: #6e56cf;
        }
        
        .swagger-ui .microlight .num {
            color: #ff9d00;
        }
        
        /* Scrollbar styling */
        ::-webkit-scrollbar {
            width: 8px;
            height: 8px;
        }
        
        ::-webkit-scrollbar-track {
            background: #1a1a24;
        }
        
        ::-webkit-scrollbar-thumb {
            background: #3d3a63;
            border-radius: 4px;
        }
        
        ::-webkit-scrollbar-thumb:hover {
            background: #6e56cf;
        }
        
        /* Markdown content */
        .swagger-ui .markdown p,
        .swagger-ui .markdown h1,
        .swagger-ui .markdown h2,
        .swagger-ui .markdown h3,
        .swagger-ui .markdown h4,
        .swagger-ui .markdown h5,
        .swagger-ui .markdown h6,
        .swagger-ui .markdown pre,
        .swagger-ui .markdown blockquote {
            color: #e0e0e8;
        }
        
        /* Authorization dialog */
        .swagger-ui .dialog-ux .modal-ux {
            background-color: #252436;
            border-radius: 8px;
        }
        
        .swagger-ui .dialog-ux .modal-ux-header h3 {
            color: #a393f0;
        }
        
        .swagger-ui .dialog-ux .modal-ux-content {
            color: #e0e0e8;
        }
        
        .swagger-ui .auth-container {
            color: #e0e0e8;
        }
        
        .swagger-ui .auth-container h4 {
            color: #a393f0;
        }
        
        .swagger-ui .scopes h2 {
            color: #a393f0;
        }
        `
