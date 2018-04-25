using System;
using Microsoft.EntityFrameworkCore;
using Dwn.Data.Database;
using Dwn.Data.Models.Identity;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.AspNetCore.Identity;
using System.Linq;

namespace Dwn.Cli.Commands
{
    public interface ICommand
    {
        void Execute();
    }
}