// angular
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
// modules
import { AlertModule } from 'ngx-bootstrap';
// services
import { AccountAPI } from '../services/AccountAPI.service'
import { BlogAPI } from '../services/BlogAPI.service'
// pages
import { BlogPage } from '../pages/Blog/Blog.page'
// components
import { AppComponent } from '../components/App/App.component';
import { BlogRollComponent } from '../components/BlogRoll/BlogRoll.component';
import { LoginComponent } from '../components/Login/Login.component';
import { PostComponent } from '../components/Post/Post.component';

@NgModule({
  declarations: [
    BlogPage,
    AppComponent,
    BlogRollComponent,
    LoginComponent,
    PostComponent
  ],
  imports: [
    AlertModule.forRoot(),
    BrowserModule
  ],
  providers: [
    AccountAPI,
    BlogAPI
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
