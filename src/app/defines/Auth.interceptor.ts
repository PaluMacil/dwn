import { Injectable } from '@angular/core';
import { HttpEvent, HttpInterceptor, HttpHandler, HttpRequest, HttpResponse, HttpErrorResponse } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import * as bs from 'bootstrap/dist/js/bootstrap.min' //required for .modal
import * as $ from 'jquery';
import 'rxjs/add/operator/map';
import 'rxjs/add/observable/throw';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = localStorage.getItem('dwn_token');
    req = req.clone({
      setHeaders: {
        Authorization: token
      }
    });
    return next.handle(req).map((event: HttpEvent<any>) => {
      return event; //do nothing
    }, (err: any) => {
      if (err instanceof HttpErrorResponse) {
        if (err.status === 401) {
          $('#login-modal').modal('show');
        } else {
          return Observable.throw(err);
        }
      }
    });
  }
}