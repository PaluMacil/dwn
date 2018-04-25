### &#8730; TODO
 - CLI Tool
   - user
     - list
     - add-role [role] [user]
   - role
     - list
     - add-user [user] [role]
 - Proxy capability (log in as any user if admin)
 - Email sender (SendGrid)
 - Expand registration form and user model (name, website, nickname)
 - SMS Two Factor (Bandwidth?)
   - Require for admin users
 - Profile picture? (after email registration)
 - Angular 6
 - Bootstrap 4
 - Font Awesome 5
 - Admin API
   - (proxy, see above)
   - User / Role Management
   - SMTP settings
   - SMTP testbed? (fake SMTP relay which acts as inbox to see all outbound emails which aren't actually sent to a real relay)
   - Menu config
   - Homepage (per role? how are multiple roles handled? also overall site setting?)

Reasons for targetting .Net Standard 2.0:
 - Quartz.NET 3
 - JSON.net 11
 - Microsoft.Data.Sqlite