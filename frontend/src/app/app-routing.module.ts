import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from './auth-guard';
import { AppComponent } from './app.component';
import { NoPermissionComponent } from './no-permission/no-permission.component';
import { HomeComponent } from './home/home.component';
import { CertComponent } from './cert/cert.component';

const routes: Routes = [{
  path: "home",
  canActivate: [AuthGuard],
  component: HomeComponent,
  data: {
    roles: ["vpn-admin-user"],
    label: "Home"
  }
},
{
  path: "cert",
  canActivate: [AuthGuard],
  component: CertComponent,
  data: {
    roles: ["vpn-admin-user", "vpn-admin-admin"],
    label: "Zertifikate"
  }
},
{
  path: "",
  canActivate: [AuthGuard],
  component: AppComponent,
  data: {
    roles: ["vpn-admin-user"]
  }
},
{
  path: "**",
  canActivate: [],
  component: NoPermissionComponent,
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
