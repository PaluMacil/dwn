// angular
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule }   from '@angular/forms';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
// modules
import { AlertModule } from 'ngx-bootstrap';
import { RoutingModule } from './Routing.module';
import { MdRenderModule } from '@nvxme/ngx-md-render';
// services
import { AccountAPI } from '../services/AccountAPI.service'
import { BlogAPI } from '../services/BlogAPI.service'
// pages
import { BlogPage } from '../pages/Blog/Blog.page'
import { BlogEditorPage } from '../pages/BlogEditor/BlogEditor.page'
// components
import { AppComponent } from '../components/App/App.component';
import { BlogRollComponent } from '../components/BlogRoll/BlogRoll.component';
import { LoginComponent } from '../components/Login/Login.component';
import { PostComponent } from '../components/Post/Post.component';
import { PostEditorComponent } from '../components/PostEditor/PostEditor.component';

@NgModule({
  declarations: [
    BlogPage,
    BlogEditorPage,

    AppComponent,
    BlogRollComponent,
    LoginComponent,
    PostComponent,
    PostEditorComponent,
  ],
  imports: [
    AlertModule.forRoot(),
    BrowserModule,
    FormsModule,
    RoutingModule,
    HttpClientModule,
    MdRenderModule,
  ],
  providers: [
    AccountAPI,
    BlogAPI
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
